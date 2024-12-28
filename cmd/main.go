package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ryoeuyo/sso/internal/app"
	"github.com/ryoeuyo/sso/internal/config"
	"github.com/ryoeuyo/sso/internal/database/postgres"
	"github.com/ryoeuyo/sso/internal/share/logger"
)

func main() {
	cfg := config.MustLoad()
	l := logger.Setup(cfg.Env)

	l.Info("Configuration loaded", slog.String("env", cfg.Env))

	repository := postgres.MustInit(cfg.Database)
	defer repository.Stop()

	l.Debug("Repository configured", slog.Any("port", cfg.Database.Port))

	application := app.New(l, cfg.GRPCServer.Port, repository, cfg.GRPCServer.TokenTTL, cfg.JWTSecretKey)
	go application.Srv.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Srv.Stop()
	l.Info("Shutting down")
}
