package app

import (
	"github.com/ryoeuyo/auth-microservice/internal/app/grpcapp"
	"github.com/ryoeuyo/auth-microservice/internal/config"
	"github.com/ryoeuyo/auth-microservice/internal/domain/entity"
	"github.com/ryoeuyo/auth-microservice/internal/service/auth"
	"log/slog"
)

type App struct {
	Srv *grpcapp.App
}

func New(
	log *slog.Logger,
	repo entity.UserRepository,
	cfg *config.AppConfig,
) *App {
	authService := auth.New(log, repo, cfg.GRPCServer.TokenTTL, cfg.JWTSecretKey)
	grpcSrv := grpcapp.New(log, authService, cfg.GRPCServer.Port)

	return &App{
		Srv: grpcSrv,
	}
}
