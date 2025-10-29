package grpc

import (
	"context"
	pb "vado_server/api/pb/ping"

	"google.golang.org/protobuf/types/known/emptypb"
)

type PingServer struct {
	pb.UnsafePingServiceServer
}

func (s *PingServer) Ping(_ context.Context, _ *emptypb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{Run: true}, nil
}
