package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"vado_server/internal/db"
	"vado_server/internal/router"
	"vado_server/internal/util"

	"time"
	"vado_server/internal/appcontext"
	"vado_server/internal/logger"

	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	pb "vado_server/internal/pb/chat"
	pbHello "vado_server/internal/pb/hello"

	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load(".env")
	appCtx := appcontext.NewAppContext(initLogger())
	appCtx.Log.Infow("Start vado-server.", "time", time.Now().Format("2006-01-02 15:04:05"))
	database := initDB(appCtx)
	appCtx.DB = database

	var wg sync.WaitGroup
	wg.Add(2)
	go startHTTPServer(appCtx, &wg, util.GetEnv("PORT"))
	go startGRPCServer(appCtx, &wg, "50051")
	wg.Wait()
}

func initLogger() *zap.SugaredLogger {
	zapLogger, zapLoggerInitErr := logger.Init(true)
	if zapLoggerInitErr != nil {
		panic(zapLoggerInitErr)
	}
	defer func() { _ = zapLogger.Sync() }()

	return zapLogger
}

func initDB(appCtx *appcontext.AppContext) *gorm.DB {
	dsn := util.GetEnv("POSTGRES_DSN")
	database, err := db.Connect(dsn)
	if err != nil {
		appCtx.Log.Fatalw("Failed to connect database", "error", err)
	}

	appCtx.Log.Infow("Connected to database", "dsn", dsn)

	return database
}

func startHTTPServer(cxt *appcontext.AppContext, wg *sync.WaitGroup, port string) {
	defer wg.Done()
	r := router.SetupRouter(cxt)

	cxt.Log.Infow("HTTP (Gin) Server starting", "port", port)
	if err := r.Run(":" + port); err != nil {
		cxt.Log.Fatalw("Server failed", "error", err)
	}
}

type grpcSrv struct {
	pbHello.UnimplementedHelloServiceServer
}

type chatServer struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[pb.ChatService_ChatStreamServer]struct{}
}

func newChatService() *chatServer {
	return &chatServer{clients: make(map[pb.ChatService_ChatStreamServer]struct{})}
}

func (s *chatServer) SendMessage(_ context.Context, msg *pb.ChatMessage) (*pb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// рассылаем всем
	for client := range s.clients {
		err := client.Send(msg)
		if err != nil {
			delete(s.clients, client)
		}
	}
	return &pb.Empty{}, nil
}

func (s *chatServer) ChatStream(empty *pb.Empty, stream pb.ChatService_ChatStreamServer) error {
	s.mu.Lock()
	s.clients[stream] = struct{}{}
	s.mu.Unlock()

	// остаёмся в стриме
	<-stream.Context().Done()
	s.mu.Lock()
	delete(s.clients, stream)
	s.mu.Unlock()
	return nil
}

func (s *grpcSrv) SeyHello(_ context.Context, req *pbHello.HelloRequest) (*pbHello.HelloResponse, error) {
	return &pbHello.HelloResponse{
		Message: fmt.Sprintf("Привет, %s! Это ответ с gRPC-сервера.", req.Name),
	}, nil
}

func startGRPCServer(appCtx *appcontext.AppContext, wg *sync.WaitGroup, port string) {
	defer wg.Done()
	lis, lisErr := net.Listen("tcp", ":"+port)
	if lisErr != nil {
		appCtx.Log.Fatalf("Error listen gRPC: %v", lisErr)
	}

	grpcServer := grpc.NewServer()
	pbHello.RegisterHelloServiceServer(grpcServer, &grpcSrv{})
	pb.RegisterChatServiceServer(grpcServer, newChatService())
	appCtx.Log.Infow("gRPC server starting", "port", port)
	if err := grpcServer.Serve(lis); err != nil {
		appCtx.Log.Fatalf("Ошибка запуска gRPC: %v", err)
	}
}
