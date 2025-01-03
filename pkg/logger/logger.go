package logger

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
	"github.com/ryoeuyo/slogdiscard"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envTest  = "test"
)

func Setup(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	// When env is local, logger is pretty
	case envLocal:
		opts := &devslog.Options{
			HandlerOptions:    &slog.HandlerOptions{Level: slog.LevelDebug},
			MaxSlicePrintSize: 10,
			SortKeys:          false,
			NewLineAfterLog:   true,
			StringerFormatter: true,
		}

		log = slog.New(devslog.NewHandler(os.Stdout, opts))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	case envTest:
		log = slog.New(slogdiscard.NewDiscardHandler())
	}

	return log
}
