package authgrpc

import (
	"context"
	"errors"
	ssov1 "github.com/ryoeuyo/mi-blog-protos/gen/go/sso"
	"github.com/ryoeuyo/sso/internal/domain/entity"
	"github.com/ryoeuyo/sso/internal/usecase/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type Server struct {
	ssov1.UnimplementedAuthServer
	l   *slog.Logger
	svc entity.AuthUseCase
}

func Register(gRPCServer *grpc.Server, authUseCase entity.AuthUseCase) {
	ssov1.RegisterAuthServer(gRPCServer, &Server{svc: authUseCase})
}

func (s *Server) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	const fn = "authgrpc.Login"

	l := s.l.With(
		slog.String("fn", fn),
		slog.String("login", req.GetLogin()),
	)

	if len(req.GetLogin()) < 8 {
		l.Warn("invalid login")

		return nil, status.Error(codes.InvalidArgument, "len login could be more than 8 symbols")
	}

	if len(req.GetPassword()) < 8 {
		l.Warn("weak password")

		return nil, status.Error(codes.InvalidArgument, "len password could be more than 8 symbols")
	}

	token, err := s.svc.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			l.Warn("invalid credentials", slog.String("error", err.Error()))

			return nil, status.Error(codes.InvalidArgument, "invalid login or password")
		}

		l.Error("failed to login user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *Server) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	const fn = "authgrpc.Login"

	l := s.l.With(
		slog.String("fn", fn),
		slog.String("login", req.GetLogin()),
	)

	if req.GetLogin() == "" {
		l.Warn("missing login")

		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if req.GetPassword() == "" {
		l.Warn("missing password")

		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	id, err := s.svc.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			l.Warn("invalid credentials", slog.String("error", err.Error()))

			return nil, status.Error(codes.InvalidArgument, "invalid login or password")
		}

		l.Error("failed register user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &ssov1.RegisterResponse{UserId: id}, nil
}

func (s *Server) IsAdmin(context.Context, *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
