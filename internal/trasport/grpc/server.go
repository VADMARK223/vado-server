package grpc

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
	pbAuth "vado_server/api/pb/auth"
	pbChat "vado_server/api/pb/chat"
	pbHello "vado_server/api/pb/hello"
	pbPing "vado_server/api/pb/ping"
	"vado_server/internal/app"
	"vado_server/internal/config/kafka/topic"
	"vado_server/internal/domain/user"
	"vado_server/internal/infra/kafka"
	"vado_server/internal/infra/persistence/gorm"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
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
			grpc.UnaryInterceptor(NewAuthInterceptor(ctx.Cfg.JwtSecret)),
		),
		listener: lis,
		log:      ctx.Log,
	}
	tokenTTL, _ := strconv.Atoi(ctx.Cfg.TokenTTL)
	refreshTTL, _ := strconv.Atoi(ctx.Cfg.RefreshTTL)
	userSvc := user.NewService(gorm.NewUserRepo(ctx), time.Duration(tokenTTL)*time.Second, time.Duration(refreshTTL)*time.Second)

	pbAuth.RegisterAuthServiceServer(s.grpcServer, NewAuthServer(userSvc, ctx.Cfg.JwtSecret))
	pbHello.RegisterHelloServiceServer(s.grpcServer, NewHelloServer(ctx.Log))
	pbPing.RegisterPingServiceServer(s.grpcServer, &PingServer{})
	producer := kafka.NewProducer(topic.ChatLog, ctx.Log, ctx.Cfg)
	pbChat.RegisterChatServiceServer(s.grpcServer, New(ctx.Log, producer))

	wrappedGrpc := grpcweb.WrapServer(
		s.grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			ctx.Log.Debugw("origin", "origin", origin)
			return true
		}),
		grpcweb.WithAllowedRequestHeaders([]string{
			"x-grpc-web", "content-type", "x-user-agent", "authorization",
		}),
	)
	portHttp := "8090"
	httpServer := &http.Server{
		Addr: ":" + portHttp,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.log.Debugw("HTTP request", "method", r.Method, "path", r.URL.Path, "headers", r.Header)
			// CORS заголовки
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"x-grpc-web, content-type, x-user-agent, authorization, accept, x-requested-with")
			w.Header().Set("Access-Control-Expose-Headers", "Grpc-Status, Grpc-Message, Grpc-Encoding, Grpc-Accept-Encoding")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("gRPC-Web server is running"))
				return
			}

			if wrappedGrpc.IsGrpcWebRequest(r) ||
				wrappedGrpc.IsAcceptableGrpcCorsRequest(r) ||
				wrappedGrpc.IsGrpcWebSocketRequest(r) {
				s.log.Debugw("gRPC-Web request", "content-type", r.Header.Get("content-type"))
				wrappedGrpc.ServeHTTP(w, r)
				return
			}
			http.NotFound(w, r)
		}),
	}

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		s.log.Infow("gRPC-Web starting", "port", portHttp)
		if errServer := httpServer.ListenAndServe(); errServer != nil && !errors.Is(errServer, http.ErrServerClosed) {
			s.log.Errorw("gRPC-Web stopped with error", "error", errServer)
		}
	}()

	return s, nil
}

func (s *Server) Start() error {
	s.log.Infow("gRPC ping starting", "address", s.listener.Addr().String())
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) GracefulStop() {
	s.log.Infow("gRPC ping graceful stopping...")
	s.grpcServer.GracefulStop()

	/*if s.httpServer != nil {
		s.log.Infow("gRPC-Web graceful stopping...")
		if err := s.httpServer.Close(); err != nil {
			s.log.Errorw("failed to close gRPC-Web server", "error", err)
		}
	}*/
}

func (s *Server) Stop() {
	s.log.Infow("gRPC ping stopping...")
	s.grpcServer.Stop()

	/*if s.httpServer != nil {
		s.log.Infow("gRPC-Web stopping...")
		if err := s.httpServer.Close(); err != nil {
			s.log.Errorw("failed to close gRPC-Web server", "error", err)
		}
	}*/
}
