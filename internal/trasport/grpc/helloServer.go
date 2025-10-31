package grpc

import (
	"context"
	"fmt"
	"vado_server/api/pb/hello"

	"go.uber.org/zap"
)

type HelloServer struct {
	hello.UnimplementedHelloServiceServer
	log *zap.SugaredLogger
}

func NewHelloServer(log *zap.SugaredLogger) *HelloServer {
	return &HelloServer{
		log: log,
	}
}

func (s *HelloServer) SayHello(_ context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	s.log.Debugw("SayHello", "name", req.Name)
	return &hello.HelloResponse{
		Message: fmt.Sprintf("Привет, %s! Твой ID в =)", req.Name),
	}, nil
}
