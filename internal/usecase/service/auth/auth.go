package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ryoeuyo/sso/internal/database"
	"github.com/ryoeuyo/sso/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	l        *slog.Logger
	Repo     entity.UserRepository
	TokenTTL time.Duration
}

func New(log *slog.Logger, repo entity.UserRepository, ttl time.Duration) *Service {
	return &Service{
		l:        log,
		Repo:     repo,
		TokenTTL: ttl,
	}
}

func (s *Service) Login(
	ctx context.Context,
	login string,
	pass string,
) (string, error) {
	const fn = "auth.Login"

	l := s.l.With(
		slog.String("fn", fn),
		slog.String("login", login),
	)

	user, err := s.Repo.User(ctx, login)
	if err != nil {
		if errors.Is(err, database.ErrUserIsNotExists) {
			l.Warn("login is not exists", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %v", fn, ErrInvalidCredentials)
		}

		l.Error("couldn't find user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass)); err != nil {
		l.Warn("invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %v", fn, ErrInvalidCredentials)
	}

	l.Info("user successfully logged")

	// TODO: token

	return "", nil
}

func (s *Service) Register(
	ctx context.Context,
	login string,
	pass string,
) (int64, error) {
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
		l.Warn("failed to save user", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}
