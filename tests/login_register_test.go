package tests

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryoeuyo/auth-microservice/pkg/testuitls"
	"github.com/ryoeuyo/auth-microservice/tests/suite"
	ssov1 "github.com/ryoeuyo/mi-blog-protos/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegister_Login_HappyPath(t *testing.T) {
	ctx, s := suite.New(t)

	login, password := testuitls.RandomLoginAndPassword(10)

	respReg, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLog, err := s.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Login:    login,
		Password: password,
	})
	require.NoError(t, err)

	token := respLog.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Cfg.JWTSecretKey), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, login, claims["login"].(string))
	assert.Equal(t, respReg.GetUserId(), int64(claims["id"].(float64)))
}

func TestRegister_FailCases(t *testing.T) {
	ctx, s := suite.New(t)

	tests := []struct {
		name        string
		login       string
		password    string
		expectedErr string
	}{
		{
			name:        "without login",
			login:       "",
			password:    "fsdfdsfsdfsd",
			expectedErr: "len login could be more than 8 symbols",
		},
		{
			name:        "without password",
			login:       "fefwffwefwe",
			password:    "",
			expectedErr: "len password could be more than 8 symbols",
		},
		{
			name:        "withoud both parametrs",
			login:       "",
			password:    "",
			expectedErr: "len login could be more than 8 symbols",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Login:    tt.login,
				Password: tt.password,
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestRegister_WithAlreadyExistsLogin(t *testing.T) {
	ctx, s := suite.New(t)
	login, password := testuitls.RandomLoginAndPassword(10)

	respReg, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "login already exists")

}
