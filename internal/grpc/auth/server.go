package auth

import (
	"context"

	sso_v1 "github.com/ViciousKit/proto/generated/go/sso"
	"google.golang.org/grpc"
)

type serverApi struct {
	sso_v1.UnimplementedAuthServer
}

func RegisterServerApi(grpc *grpc.Server) {
	sso_v1.RegisterAuthServer(grpc, &serverApi{})
}

func (s *serverApi) Login(ctx context.Context, request *sso_v1.LoginRequest) (*sso_v1.LoginResponse, error) {
	// panic("not implemented")
	return &sso_v1.LoginResponse{
		Token: "123",
	}, nil
}

func (s *serverApi) Register(ctx context.Context, request *sso_v1.RegisterRequest) (*sso_v1.RegisterResponse, error) {
	panic("not implemented")
}

func (s *serverApi) IsAdmin(ctx context.Context, request *sso_v1.IsAdminRequest) (*sso_v1.IsAdminResponse, error) {
	panic("not implemented")
}
