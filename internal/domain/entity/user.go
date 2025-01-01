package entity

import (
	"context"
)

type User struct {
	ID       int64
	Login    string
	PassHash []byte
}

type AuthService interface {
	Login(ctx context.Context, login string, pass string) (string, error)
	Register(ctx context.Context, login string, pass string) (int64, error)
}

type UserRepository interface {
	Save(ctx context.Context, login string, passHash []byte) (int64, error)
	User(ctx context.Context, login string) (*User, error)
}
