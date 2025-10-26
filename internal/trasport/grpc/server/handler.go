package server

import (
	"context"
	pb "vado_server/api/pb/server"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedServerServiceServer
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Run: true}, nil
}
