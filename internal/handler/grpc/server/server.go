package server

import (
	"context"
	pb "vado_server/api/pb/server"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ServerService struct {
	pb.UnimplementedServerServiceServer
}

func (s *ServerService) Ping(_ context.Context, _ *emptypb.Empty) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Run: true}, nil
}
