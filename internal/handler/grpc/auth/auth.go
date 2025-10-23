package auth

import (
	"context"
	"time"
	pb "vado_server/api/pb/auth"
	"vado_server/internal/appcontext"
	"vado_server/internal/auth"
	"vado_server/internal/models"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const TokenAliveMinutes = 15

type ServerGRPC struct {
	pb.UnimplementedAuthServiceServer
	AppCtx *appcontext.AppContext
}

func (s *ServerGRPC) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.AppCtx.Log.Infow("Try login", "username", req.Username)

	var user models.User
	if err := s.AppCtx.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	accessToken, err := auth.CreateToken(user.ID, []string{"user"}, TokenAliveMinutes*time.Minute)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create access token")
	}

	refreshToken, err := auth.CreateToken(user.ID, []string{"user"}, 7*24*time.Hour)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create refresh token")
	}

	// Можно сохранить refresh-токен в БД, если хочешь его инвалидацию
	// s.AppCtx.DB.Model(&user).Update("refresh_token", refreshToken)

	return &pb.LoginResponse{
		Id:           uint64(user.ID),
		Username:     user.Username,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *ServerGRPC) Refresh(_ context.Context, req *pb.RefreshRequest) (*pb.LoginResponse, error) {
	s.AppCtx.Log.Debugw("Refresh token", "refreshToken", req.RefreshToken)
	claims, errParseToken := auth.ParseToken(req.RefreshToken)
	if errParseToken != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	userID := claims.UserID

	var user models.User
	if err := s.AppCtx.DB.First(&user, userID).Error; err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	newAccess, err := auth.CreateToken(user.ID, []string{"user"}, TokenAliveMinutes*time.Minute)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create new token")
	}

	return &pb.LoginResponse{
		Id:           uint64(user.ID),
		Username:     user.Username,
		Token:        newAccess,
		RefreshToken: req.RefreshToken,
	}, nil
}
