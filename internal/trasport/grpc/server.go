package grpc

import (
	"fmt"
	"net"
	"time"
	pbAuth "vado_server/api/pb/auth"
	pbChat "vado_server/api/pb/chat"
	pbHello "vado_server/api/pb/hello"
	pbPing "vado_server/api/pb/ping"
	"vado_server/internal/app"
	"vado_server/internal/config/token"
	"vado_server/internal/domain/user"
	"vado_server/internal/infra/persistence/gorm"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	log        *zap.SugaredLogger
}

func NewServer(ctx *app.Context, port string) (*Server, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	s := &Server{
		grpcServer: grpc.NewServer(
			grpc.UnaryInterceptor(AuthInterceptor),
		),
		listener: lis,
		log:      ctx.Log,
	}
	userSvc := user.NewService(gorm.NewUserRepo(ctx), token.AccessAliveMinutes*time.Minute)

	pbAuth.RegisterAuthServiceServer(s.grpcServer, NewAuthServer(userSvc))
	pbHello.RegisterHelloServiceServer(s.grpcServer, &HelloServer{})
	pbChat.RegisterChatServiceServer(s.grpcServer, New())
	pbPing.RegisterPingServiceServer(s.grpcServer, &PingServer{})

	return s, nil
}

func (s *Server) Start() error {
	s.log.Infow("gRPC ping starting", "address", s.listener.Addr().String())
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) GracefulStop() {
	s.log.Infow("gRPC ping graceful stopping...")
	s.grpcServer.GracefulStop()
}

func (s *Server) Stop() {
	s.log.Infow("gRPC ping stopping...")
	s.grpcServer.Stop()
}
