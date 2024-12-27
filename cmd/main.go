package main

import (
	"github.com/ryoeuyo/sso/internal/app"
	"github.com/ryoeuyo/sso/internal/config"
	"github.com/ryoeuyo/sso/internal/share/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	l := logger.Setup(cfg.Env)

	l.Info("Configuration loaded", slog.Any("env", cfg.Env))

	application := app.New(l, cfg.GRPCServer.Port, nil, cfg.GRPCServer.TokenTTL)
	go application.Srv.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Srv.Stop()
	l.Info("Shutting down")
}
