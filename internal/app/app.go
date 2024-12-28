package app

import (
	grpcapp "github.com/ryoeuyo/sso/internal/app/grpc"
	"github.com/ryoeuyo/sso/internal/domain/entity"
	"github.com/ryoeuyo/sso/internal/usecase/service/auth"
	"log/slog"
	"time"
)

type App struct {
	Srv *grpcapp.App
}

func New(log *slog.Logger, port uint16, repository entity.UserRepository, tokenTTL time.Duration, JWTSecret string) *App {
	authService := auth.New(log, repository, tokenTTL, JWTSecret)
	grpcSrv := grpcapp.New(log, authService, port)

	return &App{
		Srv: grpcSrv,
	}
}
