package slogqlite_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/t0k4r/slogqlite"
)

func TestMain(t *testing.T) {
	l, err := slogqlite.New(os.Stdout, "tmp.db", &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug | slog.LevelError | slog.LevelInfo | slog.LevelWarn,
	})
	if err != nil {
		t.Fatal(err)
	}
	slog.SetDefault(slog.New(l))
	slog.Info("hello")
	slog.Warn(",")
	slog.Warn("world")
	slog.Error("!")
	slog.Debug("hellope")
	slog.Info("route", "path", "/hello/world/!")
}
