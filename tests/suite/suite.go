package suite

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"

	"github.com/ryoeuyo/auth-microservice/internal/config"
	ssov1 "github.com/ryoeuyo/mi-blog-protos/gen/go/sso"
	"google.golang.org/grpc"
)

type Suite struct {
	*testing.T
	Cfg        *config.AppConfig
	AuthClient ssov1.AuthClient
}

var (
	configPath = "../config/config-tests.yml"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoad(configPath)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPCServer.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	conn, err := grpc.NewClient(
		getTarget(&cfg.GRPCServer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create gRPC client: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(conn),
	}
}

func getTarget(cfg *config.GRPCServer) string {
	return net.JoinHostPort(cfg.Address, strconv.Itoa(int(cfg.Port)))
}
