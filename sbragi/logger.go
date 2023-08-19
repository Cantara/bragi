package sbragi

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var defaultLogger, _ = NewLogger(slog.NewTextHandler(os.Stdout, nil))

type Logger interface {
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Notice(msg string, args ...any)
	Warning(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	WithError(err error) Logger
	SetDefault()
}

type logger struct {
	handler slog.Handler
	slog    *slog.Logger
	depth   int
	ctx     context.Context
	err     error
}

func NewLogger(handler slog.Handler) (logger, error) {
	return newLogger(handler)
}

func NewDebugLogger() (logger, error) {
	return NewLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     LevelTrace,
	}))
}

func newLogger(handler slog.Handler) (logger, error) {
	return logger{
		handler: handler,
		slog:    slog.New(handler),
		ctx:     context.Background(),
	}, nil
}

func (l logger) SetDefault() {
	defaultLogger = l
}

func (l logger) Trace(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelTrace) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log(LevelTrace, msg, args...)
}

func (l logger) Debug(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelDebug) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log(LevelDebug, msg, args...)
}

func (l logger) Info(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelInfo) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log(LevelInfo, msg, args...)
}

func (l logger) Notice(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelNotice) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log(LevelNotice, msg, args...)
}

func (l logger) Warning(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelWarning) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log(LevelWarning, msg, args...)
}

func (l logger) Error(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelError) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log(LevelError, msg, args...)
}

func (l logger) Fatal(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelFatal) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log(LevelFatal, msg, args...)
	panic(msg)
}

func (l logger) WithError(err error) Logger {
	l.err = err
	//l.depth--
	return l
}

func Trace(msg string, args ...any) {
	defaultLogger.Trace(msg, args...)
}
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}
func Notice(msg string, args ...any) {
	defaultLogger.Notice(msg, args...)
}
func Warning(msg string, args ...any) {
	defaultLogger.Warning(msg, args...)
}
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}
func Fatal(msg string, args ...any) {
	defaultLogger.Fatal(msg, args...)
}
func WithError(err error) Logger {
	l := defaultLogger
	l.depth--
	return l.WithError(err)
}

// log is the low-level logging method for methods that take ...any.
// It must always be called directly by an exported logging method
// or function, because it uses a fixed call depth to obtain the pc.
func (l logger) log(level slog.Level, msg string, args ...any) {
	if !l.handler.Enabled(l.ctx, level) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3+l.depth, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	_ = l.handler.Handle(l.ctx, r)
}
