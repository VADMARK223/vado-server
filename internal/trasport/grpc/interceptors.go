package grpc

import (
	"context"
	"strings"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthInterceptor(log *zap.SugaredLogger, secret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		publicMethods := map[string]bool{
			"/ping.PingService/Ping":    true,
			"/auth.AuthService/Login":   true,
			"/auth.AuthService/Refresh": true,
		}

		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		log.Debugw("metadata", "md", md)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata отсутствует")
		}

		var token string
		// Пробуем достать токен из Authorization
		if values := md["authorization"]; len(values) > 0 {
			token = strings.TrimPrefix(values[0], "Bearer ")
			log.Debugw("token source", "type", "Authorization header")
		}

		if token == "" {
			if cookies := md["cookie"]; len(cookies) > 0 {
				for _, c := range cookies {
					for _, part := range strings.Split(c, ";") {
						part = strings.TrimSpace(part)
						if strings.HasPrefix(part, code.VadoToken+"=") {
							token = strings.TrimPrefix(part, code.VadoToken+"=")
							log.Debugw("token source", "type", "Cookie")
							break
						}
					}
				}
			}
		}

		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "token not found (no header or cookie)")
		}

		claims, err := auth.ParseToken(token, secret) // твоя функция проверки JWT
		if err != nil {
			_, _ = pp.Printf("not valid token: %v", err)
			return nil, status.Error(codes.Unauthenticated, "некорректный токен")
		}

		if claims.UserID() == 0 {
			return nil, status.Error(codes.Unauthenticated, "пустой userID в токене")
		}

		ctx = wrap(ctx, claims.UserID())

		return handler(ctx, req)
	}
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
