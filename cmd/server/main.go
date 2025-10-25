package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	pbServer "vado_server/api/pb/server"
	"vado_server/internal/auth"
	"vado_server/internal/db"
	grpcServer2 "vado_server/internal/handler/grpc/auth"
	"vado_server/internal/handler/grpc/chat"
	"vado_server/internal/handler/grpc/hello"
	"vado_server/internal/handler/grpc/server"
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

	"github.com/segmentio/kafka-go"
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
	//StartChatConsumer(util.GetEnv("KAFKA_BROKER"), util.GetEnv("KAFKA_TOPIC"), func(user, msg string) {
	StartChatConsumer(appCtx, "localhost:9094", "chat", func(key, msg string) {
		fmt.Printf("📩 Получено сообщение:\n  key=%s\n  value=%s\n", key, msg)
	})

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

func StartChatConsumer(ctx *appcontext.AppContext, broker, topic string, handle func(user string, msg string)) {
	ctx.Log.Infow("Kafka consumer starting", "broker", broker, "topic", topic)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "chat-consumer-group",
	})
	defer func(r *kafka.Reader) {
		err := r.Close()
		if err != nil {
			log.Printf("Kafka close error: %v", err)
		}
	}(r)

	fmt.Println("👂 Читаем сообщение из Kafka...")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}
		handle(string(m.Key), string(m.Value))
	}
}
