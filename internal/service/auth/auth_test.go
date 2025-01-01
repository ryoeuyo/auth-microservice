package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryoeuyo/auth-microservice/internal/domain/entity"
	"github.com/ryoeuyo/auth-microservice/internal/domain/mocks"
	"github.com/ryoeuyo/auth-microservice/internal/share/testuitls"
	"github.com/ryoeuyo/slogdiscard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestLogin_HappyPath(t *testing.T) {
	const (
		tokenTTL = time.Hour
		secret   = "secret"
	)

	ctx := context.Background()
	login, password := testuitls.RandomLoginAndPassword(10)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	testUser := &entity.User{
		ID:       53,
		Login:    login,
		PassHash: passHash,
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("User", ctx, login).Return(testUser, nil).Once()

	service := New(slogdiscard.NewDiscardLogger(), mockRepo, tokenTTL, secret)

	token, err := service.Login(ctx, login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	genTokenTime := time.Now()

	tokenParse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParse.Claims.(jwt.MapClaims)
	require.True(t, ok)

	require.Equal(t, login, claims["login"].(string))
	require.Equal(t, testUser.ID, int64(claims["id"].(float64)))

	const deltaSec = 0.2

	assert.InDelta(t, genTokenTime.Add(tokenTTL).Unix(), claims["exp"].(float64), deltaSec)
}

func TestLogin_ExpiredToken(t *testing.T) {
	const (
		tokenTTL = 1 * time.Nanosecond
		secret   = "secret"
	)

	ctx := context.Background()
	login, password := testuitls.RandomLoginAndPassword(10)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	testUser := &entity.User{
		ID:       53,
		Login:    login,
		PassHash: passHash,
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("User", ctx, login).Return(testUser, nil).Once()

	service := New(slogdiscard.NewDiscardLogger(), mockRepo, tokenTTL, secret)

	token, err := service.Login(ctx, login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	time.Sleep(tokenTTL)

	_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.Error(t, err)
}

func TestLogin_InvalidPassword(t *testing.T) {
	const (
		tokenTTL = 1 * time.Nanosecond
		secret   = "secret"
	)

	ctx := context.Background()
	login, password := testuitls.RandomLoginAndPassword(10)

	passHash, err := bcrypt.GenerateFromPassword([]byte("test_invalid_pass"), bcrypt.DefaultCost)
	require.NoError(t, err)

	testUser := &entity.User{
		ID:       53,
		Login:    login,
		PassHash: passHash,
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("User", ctx, login).Return(testUser, nil).Once()

	service := New(slogdiscard.NewDiscardLogger(), mockRepo, tokenTTL, secret)

	_, err = service.Login(ctx, login, password)
	require.Error(t, err)
}
