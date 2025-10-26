package hello

import (
	"context"
	"fmt"
	"vado_server/api/pb/hello"
)

type Server struct {
	hello.UnimplementedHelloServiceServer
}

func (s *Server) SeyHello(_ context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{
		Message: fmt.Sprintf("Привет, %s! Твой ID в =)", req.Name),
	}, nil
}
