package grpc

import (
	"fmt"
	"net"
	"time"
	pbAuth "vado_server/api/pb/auth"
	pbChat "vado_server/api/pb/chat"
	pbHello "vado_server/api/pb/hello"
	pbPing "vado_server/api/pb/ping"
	"vado_server/internal/app/context"
	"vado_server/internal/config/token"
	"vado_server/internal/domain/user"
	"vado_server/internal/infra/persistence/gorm"

	"vado_server/internal/trasport/grpc/chat"
	"vado_server/internal/trasport/grpc/hello"
	"vado_server/internal/trasport/grpc/ping"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	log        *zap.SugaredLogger
}

func NewServer(ctx *context.AppContext, port string) (*Server, error) {
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
	authServer := NewAuthServer(userSvc)

	pbAuth.RegisterAuthServiceServer(s.grpcServer, authServer)
	pbHello.RegisterHelloServiceServer(s.grpcServer, &hello.Server{})
	pbChat.RegisterChatServiceServer(s.grpcServer, chat.New())
	pbPing.RegisterPingServiceServer(s.grpcServer, &ping.Server{})

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
