package grpc

import (
	"context"
	"fmt"
	"vado_server/api/pb/hello"
)

type HelloServer struct {
	hello.UnimplementedHelloServiceServer
}

func (s *HelloServer) SayHello(_ context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{
		Message: fmt.Sprintf("Привет, %s! Твой ID в =)", req.Name),
	}, nil
}
