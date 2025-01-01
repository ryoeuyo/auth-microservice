package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ryoeuyo/auth-microservice/internal/database"
	"github.com/ryoeuyo/auth-microservice/internal/domain/entity"
	"github.com/ryoeuyo/auth-microservice/internal/share/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	l         *slog.Logger
	Repo      entity.UserRepository
	JWTSecret string
	TokenTTL  time.Duration
}

func New(log *slog.Logger, repo entity.UserRepository, ttl time.Duration, JWTSecret string) *Service {
	return &Service{
		l:         log,
		Repo:      repo,
		TokenTTL:  ttl,
		JWTSecret: JWTSecret,
	}
}

func (s *Service) Login(ctx context.Context, login string, pass string) (string, error) {
	const fn = "auth.Login"

	l := s.l.With(
		slog.String("fn", fn),
		slog.String("login", login),
	)

	user, err := s.Repo.User(ctx, login)
	if err != nil {
		if errors.Is(err, database.ErrUserIsNotExists) {
			l.Warn("login is not exists", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %v", fn, ErrUserNotFound)
		}

		l.Warn("couldn't find user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass)); err != nil {
		l.Warn("invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", fn, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, s.TokenTTL, s.JWTSecret)
	if err != nil {
		l.Error("failed to generate jwt token", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return token, nil
}

func (s *Service) Register(ctx context.Context, login string, pass string) (int64, error) {
	const fn = "auth.Register"

	l := s.l.With(
		slog.String("fn", fn),
		slog.String("login", login),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		l.Error("failed to generate password hash", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	id, err := s.Repo.Save(ctx, login, passHash)
	if err != nil {
		if errors.Is(err, database.ErrLoginIsExists) {
			l.Warn("login already exists", slog.String("error", err.Error()))

			return 0, fmt.Errorf("%s: %w", fn, ErrUserIsExists)
		}

		l.Warn("failed to save user", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}
