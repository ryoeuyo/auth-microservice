package main

import (
	"github.com/ryoeuyo/auth-microservice/internal/database/postgres"
	"github.com/ryoeuyo/auth-microservice/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ryoeuyo/auth-microservice/internal/app"
	"github.com/ryoeuyo/auth-microservice/internal/config"
)

func main() {
	cfg := config.MustLoad()
	l := logger.Setup(cfg.Env)

	l.Info("Config loaded", slog.String("env", cfg.Env))

	repository := postgres.MustInit(&cfg.Database)
	defer repository.Stop()

	application := app.New(l, repository, cfg)
	go application.Srv.MustStart()
	go application.MetricServer.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Srv.Stop()
	l.Info("Shutting down")
}
