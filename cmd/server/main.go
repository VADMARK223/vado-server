package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	pbServer "vado_server/api/pb/server"
	"vado_server/internal/auth"
	"vado_server/internal/db"
	grpcServer2 "vado_server/internal/handler/grpc/auth"
	"vado_server/internal/handler/grpc/chat"
	"vado_server/internal/handler/grpc/hello"
	"vado_server/internal/handler/grpc/server"
	kafka2 "vado_server/internal/kafka"
	"vado_server/internal/router"
	"vado_server/internal/util"

	"time"
	"vado_server/internal/appcontext"
	"vado_server/internal/logger"

	"github.com/k0kubun/pp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	pbAuth "vado_server/api/pb/auth"
	pb "vado_server/api/pb/chat"
	pbHello "vado_server/api/pb/hello"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load(".env")

	zapLogger := logger.Init(true)
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

	// Kafka
	consumer := kafka2.NewConsumer(appCtx)
	go func() {
		consumerRun := consumer.Run(ctx, func(key, value []byte) error {
			user := string(key)
			msg := string(value)
			appCtx.Log.Infow("Processing message", "user", user, "msg", msg)
			return nil
		})

		if consumerRun != nil {
			appCtx.Log.Errorw("Consumer stopped", "error", consumerRun)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	appCtx.Log.Info("Shutdown signal received")
	// даём сервисам небольшое время завершиться
	cancel()
	_ = consumer.Close()

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
	)

	// регистрируем сервисы
	pbAuth.RegisterAuthServiceServer(grpcServerInstance, &grpcServer2.ServerGRPC{AppCtx: appCtx})
	pbHello.RegisterHelloServiceServer(grpcServerInstance, &hello.Server{})
	pb.RegisterChatServiceServer(grpcServerInstance, chat.NewChatService())
	pbServer.RegisterServerServiceServer(grpcServerInstance, &server.Service{})

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
	if strings.Contains(info.FullMethod, "Refresh") {
		return handler(ctx, req)
	}

	_, _ = pp.Println(info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata отсутствует")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "токен не найден")
	}

	token := strings.TrimPrefix(values[0], "Bearer ")
	claims, err := auth.ParseToken(token) // твоя функция проверки JWT
	if err != nil {
		_, _ = pp.Printf("not valid token: %v", err)
		return nil, status.Error(codes.Unauthenticated, "некорректный токен")
	}

	if claims.UserID == 0 {
		return nil, status.Error(codes.Unauthenticated, "пустой userID в токене")
	}

	ctx = auth.Wrap(ctx, claims.UserID)

	return handler(ctx, req)
}
