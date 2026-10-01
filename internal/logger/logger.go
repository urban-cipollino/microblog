package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	New,
)

type Env string

const (
	Dev  Env = "dev"
	Prod Env = "prod"
	Test Env = "test"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func New(mode Env) Logger {
	switch mode {
	case Dev:
		if err := os.MkdirAll("logs", 0o750); err != nil {
			panic(err)
		}

		file, err := os.OpenFile("logs/dev.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		if err != nil {
			panic(err)
		}

		return slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case Prod:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	case Test:
		return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
}
