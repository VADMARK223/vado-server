package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"vado_server/internal/constants/code"
	"vado_server/internal/db"
	grpcServer2 "vado_server/internal/grpc/auth"
	"vado_server/internal/grpc/server"
	"vado_server/internal/middleware"
	pbServer "vado_server/internal/pb/server"
	"vado_server/internal/router"
	"vado_server/internal/util"

	"time"
	"vado_server/internal/appcontext"
	"vado_server/internal/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	pbAuth "vado_server/internal/pb/auth"
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

type helloServer struct {
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

func (s *chatServer) ChatStream(_ *pb.Empty, stream pb.ChatService_ChatStreamServer) error {
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

func (s *helloServer) SeyHello(ctx context.Context, req *pbHello.HelloRequest) (*pbHello.HelloResponse, error) {
	userId := ctx.Value(code.UserId)
	return &pbHello.HelloResponse{
		Message: fmt.Sprintf("Привет, %s! Твой ID в БД=%f", req.Name, userId.(float64)),
	}, nil
}

func startGRPCServer(appCtx *appcontext.AppContext, wg *sync.WaitGroup, port string) {
	defer wg.Done()
	lis, lisErr := net.Listen("tcp", ":"+port)
	if lisErr != nil {
		appCtx.Log.Fatalf("Error listen gRPC: %v", lisErr)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor), // Перехват для обычных запросов
		//grpc.UnaryInterceptor(AuthStreamInterceptor),
	)
	pbAuth.RegisterAuthServiceServer(grpcServer, &grpcServer2.AuthServerGRPC{AppCtx: appCtx})
	pbHello.RegisterHelloServiceServer(grpcServer, &helloServer{})
	pb.RegisterChatServiceServer(grpcServer, newChatService())
	pbServer.RegisterServerServiceServer(grpcServer, &server.ServerService{})
	appCtx.Log.Infow("gRPC server starting", "port", port)
	if err := grpcServer.Serve(lis); err != nil {
		appCtx.Log.Fatalf("Ошибка запуска gRPC: %v", err)
	}
}

func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Не проверяем токен для публичных методов
	if strings.Contains(info.FullMethod, "Ping") {
		return handler(ctx, req)
	}
	if strings.Contains(info.FullMethod, "Login") {
		return handler(ctx, req)
	}

	// Достаём токен из metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata отсутствует")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "токен не найден")
	}

	token := strings.TrimPrefix(values[0], "Bearer ")
	claims, err := middleware.ParseToken(token) // твоя функция проверки JWT
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "некорректный токен")
	}

	if claims.UserID != 0 {
		ctx = context.WithValue(ctx, code.UserId, claims.UserID)
	}

	return handler(ctx, req)
}

/*func AuthStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	if strings.Contains(info.FullMethod, "Login") {
		return handler(srv, ss)
	}

	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata отсутствует")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Error(codes.Unauthenticated, "токен не найден")
	}

	token := strings.TrimPrefix(values[0], "Bearer ")
	claims, err := auth.ParseToken(token)
	if err != nil {
		return status.Error(codes.Unauthenticated, "токен невалиден")
	}

	// Оборачиваем stream с контекстом, где уже есть userID
	wrapped := &wrappedStream{
		ServerStream: ss,
		ctx:          context.WithValue(ss.Context(), code.UserId, claims.UserID),
	}

	return handler(srv, wrapped)
}

// вспомогательная обёртка
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}*/
