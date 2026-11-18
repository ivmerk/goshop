//go:build integration

package integration

import (
	"context"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authV1 "github.com/ivmerk/goshop/shared/pkg/proto/auth/v1"
	commonV1 "github.com/ivmerk/goshop/shared/pkg/proto/common/v1"
)

// fakeIAM заглушка AuthService: признаёт валидной только одну сессию.
type fakeIAM struct {
	authV1.UnimplementedAuthServiceServer

	validSession string
	server       *grpc.Server
	port         int
}

func startFakeIAM(validSession string) (*fakeIAM, error) {
	// слушаем на всех интерфейсах: контейнер ходит к хосту через host-gateway
	lis, err := net.Listen("tcp", "0.0.0.0:0") //nolint:gosec // тестовая заглушка, должна быть доступна из docker
	if err != nil {
		return nil, err
	}

	f := &fakeIAM{
		validSession: validSession,
		server:       grpc.NewServer(),
		port:         lis.Addr().(*net.TCPAddr).Port,
	}
	authV1.RegisterAuthServiceServer(f.server, f)

	go func() { _ = f.server.Serve(lis) }() //nolint:errcheck,gosec // ошибка Serve после Stop не важна

	return f, nil
}

func (f *fakeIAM) Port() string { return strconv.Itoa(f.port) }

func (f *fakeIAM) Stop() { f.server.GracefulStop() }

func (f *fakeIAM) Whoami(_ context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	if req.GetSessionUuid() != f.validSession {
		return nil, status.Error(codes.Unauthenticated, "session not found")
	}

	return &authV1.WhoamiResponse{UserInfo: &commonV1.UserInfo{
		UserUuid: "22222222-2222-2222-2222-222222222222",
		Login:    "test-user",
		Email:    "test@example.com",
	}}, nil
}
