package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryoeuyo/auth-microservice/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	secret   = "test-secret"
	tokenTTL = 15 * time.Minute
)

func TestNewToken_HappyPath(t *testing.T) {
	testUser := &entity.User{
		ID:       99,
		Login:    "testUser",
		PassHash: []byte("testPass"),
	}

	token, err := NewToken(testUser, tokenTTL, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	genTokenTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, testUser.Login, claims["login"].(string))
	assert.Equal(t, testUser.ID, int64(claims["id"].(float64)))

	const deltaSec = 0.2

	assert.InDelta(t, genTokenTime.Add(tokenTTL).Unix(), claims["exp"].(float64), deltaSec)
}

func TestNewToken_Expired(t *testing.T) {
	testUser := &entity.User{
		ID:       99,
		Login:    "testUser",
		PassHash: []byte("testPass"),
	}

	smallTTL := 1 * time.Nanosecond

	token, err := NewToken(testUser, smallTTL, secret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	time.Sleep(smallTTL)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.Error(t, err)
	assert.Equal(t, tokenParsed.Valid, false)
}

func TestNewToken_InvalidSecret(t *testing.T) {
	testUser := &entity.User{
		ID:       99,
		Login:    "testUser",
		PassHash: []byte("testPass"),
	}

	token, err := NewToken(testUser, tokenTTL, secret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("invalid-secret"), nil
	})
	require.Error(t, err)
	assert.Equal(t, tokenParsed.Valid, false)
}
