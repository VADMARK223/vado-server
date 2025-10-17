package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	pbServer "vado_server/api/pb/server"
	"vado_server/internal/constants/code"
	"vado_server/internal/db"
	grpcServer2 "vado_server/internal/handler/grpc/auth"
	"vado_server/internal/handler/grpc/chat"
	"vado_server/internal/handler/grpc/hello"
	"vado_server/internal/handler/grpc/server"
	"vado_server/internal/middleware"
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

	pbAuth "vado_server/api/pb/auth"
	pb "vado_server/api/pb/chat"
	pbHello "vado_server/api/pb/hello"

	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load(".env")

	zapLogger := initLogger()
	defer func() { _ = zapLogger.Sync() }()

	appCtx := appcontext.NewAppContext(zapLogger)
	appCtx.Log.Infow("Start vado-server.", "time", time.Now().Format("2006-01-02 15:04:05"))

	database := initDB(appCtx)
	appCtx.DB = database
	defer func() {
		if sqlDB, err := database.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	// HTTP сервер
	wg.Add(1)
	go startHTTPServer(ctx, appCtx, &wg, util.GetEnv("PORT"))

	// gRPC сервер
	wg.Add(1)
	grpcServerInstance, err := startGRPCServer(appCtx, &wg, util.GetEnv("GRPC_PORT"))
	if err != nil {
		appCtx.Log.Fatalw("failed to start grpc server", "error", err)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	appCtx.Log.Info("Shutdown signal received")
	// даём сервисам небольшое время завершиться
	cancel()

	// Останавливаем gRPC сервер корректно
	if grpcServerInstance != nil {
		// GracefulStop не принимает контекст; оборачиваем в go func чтобы не блокировать
		done := make(chan struct{})
		go func() {
			appCtx.Log.Info("gRPC: GracefulStop called")
			grpcServerInstance.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			appCtx.Log.Info("gRPC server stopped gracefully")
		case <-time.After(10 * time.Second):
			appCtx.Log.Warn("gRPC graceful stop timeout, forcing Stop()")
			grpcServerInstance.Stop()
		}
	}

	wg.Wait()
	appCtx.Log.Infow("Servers stopped.")
}

func initLogger() *zap.SugaredLogger {
	zapLogger, zapLoggerInitErr := logger.Init(true)
	if zapLoggerInitErr != nil {
		panic(zapLoggerInitErr)
	}

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

func startHTTPServer(ctx context.Context, appCtx *appcontext.AppContext, wg *sync.WaitGroup, port string) {
	defer wg.Done()

	r := router.SetupRouter(appCtx)
	appCtx.Log.Infow("HTTP (Gin) Server starting", "port", port)

	// Запускаем сервер в отдельной горутине для graceful shutdown
	go func() {
		if err := r.Run(":" + port); err != nil {
			appCtx.Log.Errorw("HTTP Server error", "error", err)
		}
	}()

	<-ctx.Done()
	appCtx.Log.Info("HTTP Server shutting down")
}

func startGRPCServer(appCtx *appcontext.AppContext, wg *sync.WaitGroup, port string) (*grpc.Server, error) {
	lis, lisErr := net.Listen("tcp", ":"+port)
	if lisErr != nil {
		return nil, lisErr
	}

	grpcServerInstance := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor), // Перехват для обычных запросов
		//grpc.UnaryInterceptor(AuthStreamInterceptor),
	)

	// регистрируем сервисы
	pbAuth.RegisterAuthServiceServer(grpcServerInstance, &grpcServer2.AuthServerGRPC{AppCtx: appCtx})
	pbHello.RegisterHelloServiceServer(grpcServerInstance, &hello.HelloServer{})
	pb.RegisterChatServiceServer(grpcServerInstance, chat.NewChatService())
	pbServer.RegisterServerServiceServer(grpcServerInstance, &server.ServerService{})

	appCtx.Log.Infow("gRPC server starting", "port", port)

	go func() {
		defer wg.Done()
		if err := grpcServerInstance.Serve(lis); err != nil {
			appCtx.Log.Errorw("gRPC Server Serve error", "error", err)
		}
		appCtx.Log.Info("gRPC Serve returned")
	}()

	// возвращаем сервер сразу (негатив: Serve находится в горутине; остановка — через GracefulStop в main)
	return grpcServerInstance, nil
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
