package auth

import (
	"context"
	"time"
	"vado_server/internal/appcontext"
	"vado_server/internal/auth"
	"vado_server/internal/models"
	pb "vado_server/internal/pb/auth"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServerGRPC struct {
	pb.UnimplementedAuthServiceServer
	AppCtx *appcontext.AppContext
}

func (s *AuthServerGRPC) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.AppCtx.Log.Infow("LOGIN", "username", req.Username, "password", req.Password)

	var user models.User
	if err := s.AppCtx.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	token, err := auth.CreateToken(user.ID, []string{"user"}, 15*time.Minute)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create token")
	}

	return &pb.LoginResponse{Token: token}, nil
}
