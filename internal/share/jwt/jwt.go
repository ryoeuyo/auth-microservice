package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryoeuyo/sso/internal/domain/entity"
)

func NewToken(user *entity.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
