package auth

import (
	"context"

	sso_v1 "github.com/ViciousKit/proto/generated/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appId int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int) (bool, error)
}

type serverApi struct {
	sso_v1.UnimplementedAuthServer
	auth Auth
}

func RegisterServerApi(grpc *grpc.Server, auth Auth) {
	sso_v1.RegisterAuthServer(grpc, &serverApi{auth: auth})
}

func (s *serverApi) Login(ctx context.Context, req *sso_v1.LoginRequest) (*sso_v1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Password required")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Appid required")
	}

	//TODO: login and return token
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso_v1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverApi) Register(ctx context.Context, req *sso_v1.RegisterRequest) (*sso_v1.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Password required")
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso_v1.RegisterResponse{
		UserId: userId,
	}, nil
}

func (s *serverApi) IsAdmin(ctx context.Context, req *sso_v1.IsAdminRequest) (*sso_v1.IsAdminResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "UserId required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, int(req.GetUserId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso_v1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
