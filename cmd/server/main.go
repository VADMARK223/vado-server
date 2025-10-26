package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	ctx "vado_server/internal/app/context"
	"vado_server/internal/app/logger"
	"vado_server/internal/infra/db"
	"vado_server/internal/infra/kafka"
	"vado_server/internal/trasport/grpc"
	"vado_server/internal/trasport/http"
	"vado_server/internal/util"

	"time"

	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	zapLogger := logger.Init(true)
	defer func() { _ = zapLogger.Sync() }()

	appCtx := ctx.NewAppContext(zapLogger)
	appCtx.Log.Infow("Start vado-server.", "time", time.Now().Format("2006-01-02 15:04:05"))

	database := initDB(appCtx)
	appCtx.DB = database
	defer func() {
		if sqlDB, err := database.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}()

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	// HTTP сервер
	wg.Add(1)
	go startHTTPServer(ctxWithCancel, appCtx, &wg, util.GetEnv("PORT"))

	// gRPC сервер
	grpcServer, err := grpc.NewServer(appCtx, util.GetEnv("GRPC_PORT"))
	if err != nil {
		appCtx.Log.Fatalw("failed to start grpc server", "error", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Start(); err != nil {
			appCtx.Log.Errorw("grpc server stopped", "error", err)
		}
	}()

	// Kafka
	consumer := kafka.NewConsumer(appCtx)
	go func() {
		consumerRun := consumer.Run(ctxWithCancel, func(key, value []byte) error {
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
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	appCtx.Log.Info("Shutdown signal received")
	// даём сервисам небольшое время завершиться
	cancel()
	_ = consumer.Close()

	// Останавливаем gRPC сервер корректно
	if grpcServer != nil {
		// GracefulStop не принимает контекст; оборачиваем в горутину, чтобы не блокировать
		done := make(chan struct{})
		go func() {
			appCtx.Log.Info("gRPC: GracefulStop called")
			grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			appCtx.Log.Info("gRPC server stopped gracefully")
		case <-time.After(10 * time.Second):
			appCtx.Log.Warn("gRPC graceful stop timeout, forcing Stop()")
			grpcServer.Stop()
		}
	}

	wg.Wait()
	appCtx.Log.Infow("Servers stopped.")
}

func initDB(appCtx *ctx.AppContext) *gorm.DB {
	dsn := util.GetEnv("POSTGRES_DSN")
	database, err := db.Connect(dsn)
	if err != nil {
		appCtx.Log.Fatalw("Failed to connect database", "error", err)
	}

	appCtx.Log.Infow("Connected to database", "dsn", dsn)

	return database
}

func startHTTPServer(ctx context.Context, appCtx *ctx.AppContext, wg *sync.WaitGroup, port string) {
	defer wg.Done()

	r := http.SetupRouter(appCtx)
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
