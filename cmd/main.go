package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ryoeuyo/sso/internal/app"
	"github.com/ryoeuyo/sso/internal/config"
	"github.com/ryoeuyo/sso/internal/share/logger"
)

func main() {
	cfg := config.MustLoad()
	l := logger.Setup(cfg.Env)

	l.Info("Config loaded", slog.String("env", cfg.Env))

	application := app.New(l, cfg)
	go application.Srv.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Srv.Stop()
	l.Info("Shutting down")
}
