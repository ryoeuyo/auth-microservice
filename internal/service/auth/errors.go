package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserIsExists       = errors.New("user already exists")
)
