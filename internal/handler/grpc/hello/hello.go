package hello

import (
	"context"
	"fmt"
	"vado_server/api/pb/hello"
	"vado_server/internal/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	hello.UnimplementedHelloServiceServer
}

func (s *Server) SeyHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	userId, ok := auth.TryGet(ctx)

	if !ok {
		return nil, status.Error(codes.Unauthenticated, "userID не найден")
	}

	return &hello.HelloResponse{
		Message: fmt.Sprintf("Привет, %s! Твой ID в БД=%d =)", req.Name, userId),
	}, nil
}
