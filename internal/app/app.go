package app

import (
	"github.com/ryoeuyo/sso/internal/app/grpcapp"
	"github.com/ryoeuyo/sso/internal/config"
	"github.com/ryoeuyo/sso/internal/database/postgres"
	"github.com/ryoeuyo/sso/internal/usecase/service/auth"
	"log/slog"
)

type App struct {
	Srv *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.AppConfig,
) *App {
	repository := postgres.MustInit(&cfg.Database)
	defer repository.Stop()

	authService := auth.New(log, repository, cfg.GRPCServer.TokenTTL, cfg.JWTSecretKey)
	grpcSrv := grpcapp.New(log, authService, cfg.GRPCServer.Port)

	return &App{
		Srv: grpcSrv,
	}
}
