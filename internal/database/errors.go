package database

import "errors"

var (
	ErrLoginIsExists   = errors.New("login is exists")
	ErrUserIsNotExists = errors.New("user is not exists")
)
