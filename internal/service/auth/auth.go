package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/ryoeuyo/auth-microservice/internal/app/metric"
	"github.com/ryoeuyo/auth-microservice/pkg/jwt"
	"log/slog"
	"time"

	"github.com/ryoeuyo/auth-microservice/internal/database"
	"github.com/ryoeuyo/auth-microservice/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	l         *slog.Logger
	Repo      entity.UserRepository
	Metric    *metric.Metric
	JWTSecret string
	TokenTTL  time.Duration
}

func New(log *slog.Logger, repo entity.UserRepository, metric *metric.Metric, ttl time.Duration, JWTSecret string) *Service {
	return &Service{
		l:         log,
		Repo:      repo,
		Metric:    metric,
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

	mFail := s.Metric.AuthFailedAttempts.WithLabelValues("login")
	mReqCounter := s.Metric.AuthRequests.WithLabelValues("login")
	mReqDuration := s.Metric.AuthRequestDuration.WithLabelValues("login")
	start := time.Now()

	errCh := make(chan error, 1)

	defer func(start time.Time) {
		mReqDuration.Observe(time.Since(start).Seconds())

		if err := <-errCh; err != nil {
			mFail.Inc()
		}

		mReqCounter.Inc()
	}(start)

	user, err := s.Repo.User(ctx, login)
	if err != nil {
		errCh <- err
		if errors.Is(err, database.ErrUserIsNotExists) {
			l.Warn("login is not exists", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %v", fn, ErrUserNotFound)
		}
		l.Warn("couldn't find user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass)); err != nil {
		errCh <- err
		l.Warn("invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", fn, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, s.TokenTTL, s.JWTSecret)
	if err != nil {
		errCh <- err
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

	mFail := s.Metric.AuthFailedAttempts.WithLabelValues("register")
	mReqCounter := s.Metric.AuthRequests.WithLabelValues("register")
	mReqDuration := s.Metric.AuthRequestDuration.WithLabelValues("register")
	start := time.Now()

	errCh := make(chan error, 1)

	defer func(start time.Time) {
		mReqDuration.Observe(time.Since(start).Seconds())

		if err := <-errCh; err != nil {
			mFail.Inc()
		}

		mReqCounter.Inc()
	}(start)

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		errCh <- err
		l.Error("failed to generate password hash", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	id, err := s.Repo.Save(ctx, login, passHash)
	if err != nil {
		errCh <- err
		if errors.Is(err, database.ErrLoginIsExists) {
			l.Warn("login already exists", slog.String("error", err.Error()))

			return 0, fmt.Errorf("%s: %w", fn, ErrUserIsExists)
		}

		mFail.Inc()
		l.Warn("failed to save user", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}
