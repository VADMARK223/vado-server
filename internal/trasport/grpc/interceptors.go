package grpc

import (
	"context"
	"strings"
	"vado_server/internal/auth"

	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Не проверяем токен для публичных методов
	if strings.Contains(info.FullMethod, "Ping") {
		return handler(ctx, req)
	}
	if strings.Contains(info.FullMethod, "Login") {
		return handler(ctx, req)
	}
	if strings.Contains(info.FullMethod, "Refresh") {
		return handler(ctx, req)
	}

	_, _ = pp.Println(info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata отсутствует")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "токен не найден")
	}

	token := strings.TrimPrefix(values[0], "Bearer ")
	claims, err := auth.ParseToken(token) // твоя функция проверки JWT
	if err != nil {
		_, _ = pp.Printf("not valid token: %v", err)
		return nil, status.Error(codes.Unauthenticated, "некорректный токен")
	}

	if claims.UserID == 0 {
		return nil, status.Error(codes.Unauthenticated, "пустой userID в токене")
	}

	ctx = wrap(ctx, claims.UserID)

	return handler(ctx, req)
}

type AuthContext struct {
	context.Context
	userID uint
}

func (a *AuthContext) UserID() uint {
	return a.userID
}

func wrap(ctx context.Context, userID uint) context.Context {
	return &AuthContext{
		Context: ctx,
		userID:  userID,
	}
}
