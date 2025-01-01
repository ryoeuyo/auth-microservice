package testuitls

import (
	"math/rand"
	"time"
)

func RandomLoginAndPassword(length int) (string, string) {
	const (
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	)

	rand.New(rand.NewSource(time.Now().UnixNano()))

	login := make([]byte, length)
	for i := range login {
		login[i] = charset[rand.Intn(len(charset))]
	}

	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}

	return string(login), string(password)
}
