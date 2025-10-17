package hello

import (
	"context"
	"fmt"
	"vado_server/api/pb/hello"
	"vado_server/internal/constants/code"
)

type Server struct {
	hello.UnimplementedHelloServiceServer
}

func (s *Server) SeyHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	userId := ctx.Value(code.UserId)
	return &hello.HelloResponse{
		Message: fmt.Sprintf("Привет, %s! Твой ID в БД=%d =)", req.Name, userId),
	}, nil
}
