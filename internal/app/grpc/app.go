package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/ryoeuyo/sso/internal/domain/entity"
	"github.com/ryoeuyo/sso/internal/transport/grpc/authgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type App struct {
	log    *slog.Logger
	Server *grpc.Server
	port   uint16
}

func New(log *slog.Logger, useCase entity.AuthUseCase, port uint16) *App {
	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
	}

	recoverOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoverOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), logOpts...),
	))

	authgrpc.Register(server, useCase)

	return &App{
		log:    log,
		Server: server,
		port:   port,
	}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) MustStart() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		panic(err)
	}

	a.log.Info("grpc server listening on address", slog.String("address", l.Addr().String()))

	if err := a.Server.Serve(l); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const fn = "grpc.Stop"

	a.log.With(slog.String("fn", fn)).
		Info("stopping grpc server", slog.Any("port", a.port))

	a.Server.GracefulStop()
}
