package grpc

import (
	"context"
	pb "vado_server/api/pb/auth"
	"vado_server/internal/domain/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	service *user.Service
}

func NewAuthServer(service *user.Service) *AuthServer {
	return &AuthServer{
		service: service,
	}
}

func (s *AuthServer) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	username := req.Username
	password := req.Password

	u, accessToken, refreshToken, err := s.service.Login(username, password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.LoginResponse{
		Id:           uint64(u.ID),
		Username:     u.Username,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServer) Refresh(_ context.Context, req *pb.RefreshRequest) (*pb.LoginResponse, error) {
	u, newToken, err := s.service.Refresh(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if u == nil {
		return nil, status.Error(codes.Internal, "user not found")
	}

	return &pb.LoginResponse{
		Id:           uint64(u.ID),
		Username:     u.Username,
		Token:        newToken,
		RefreshToken: req.RefreshToken,
	}, nil
}
