package sbragi

import (
	"golang.org/x/exp/slog"
	"testing"
	"time"
)

func TestLongRunning(t *testing.T) {
	h, err := NewHandlerInFolder("./log")
	if err != nil {
		t.Error(err)
		return
	}
	defer h.Cancel()
	h.level = slog.LevelDebug
	log := slog.New(&h)
	i := 0
	ticker := time.Tick(time.Nanosecond)
	for range ticker {
		if i%2 == 0 {
			log.Log(slog.LevelInfo, "text", "error", i)
		} else {
			log.Debug("text", "error", i)
		}
		i++
		if i > 10000 {
			break
		}
	}
}
