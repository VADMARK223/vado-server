package ping

import (
	"context"
	pb "vado_server/api/pb/ping"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnsafePingServiceServer
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{Run: true}, nil
}
