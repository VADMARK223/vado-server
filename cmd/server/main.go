package main

import (
	"context"
	"errors"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	ctx "vado_server/internal/app"
	"vado_server/internal/config/code"
	"vado_server/internal/infra/db"
	"vado_server/internal/infra/kafka"
	"vado_server/internal/infra/logger"
	"vado_server/internal/trasport/grpc"
	"vado_server/internal/trasport/http"

	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func main() {
	//------------------------------------------------------------
	// Загрузка окружения
	//------------------------------------------------------------
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = code.Local // по умолчанию, если не задано
	}
	switch env {
	case code.Local:
		if err := godotenv.Load(".env.local"); err != nil {
			log.Println("⚠️  .env.local not found — using system env")
		} else {
			log.Println("✅ Loaded .env.local")
		}
	default:
		log.Println("ℹ️  Running in", env, "mode — skipping local env")
	}
	//------------------------------------------------------------
	// Инициализация логгера и контекста приложения
	//------------------------------------------------------------
	zapLogger := logger.Init(true)
	defer func() { _ = zapLogger.Sync() }()

	appCtx := ctx.NewAppContext(zapLogger)
	appCtx.Log.Infow("Start vado-ping.", "time", time.Now().Format("2006-01-02 15:04:05"))

	//------------------------------------------------------------
	// Подключение к базе данных
	//------------------------------------------------------------
	database := initDB(appCtx)
	appCtx.DB = database
	defer func() {
		if sqlDB, err := database.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}()

	//------------------------------------------------------------
	// Общий контекст и группа ожидания
	//------------------------------------------------------------
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	defer cancel()

	//------------------------------------------------------------
	// HTTP сервер (Gin)
	//------------------------------------------------------------
	wg.Add(1)
	go startHTTPServer(ctxWithCancel, appCtx, &wg, appCtx.Cfg.Port)

	//------------------------------------------------------------
	// gRPC сервер
	//------------------------------------------------------------
	grpcServer, err := grpc.NewServer(appCtx, appCtx.Cfg.GrpcPort)
	if err != nil {
		appCtx.Log.Fatalw("failed to start gRPC server", "error", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Start(); err != nil {
			appCtx.Log.Errorw("gRPC server stopped", "error", err)
		}
	}()

	//------------------------------------------------------------
	// Kafka consumer
	//------------------------------------------------------------
	consumer := kafka.NewConsumer(appCtx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr := consumer.Run(ctxWithCancel, func(key, value []byte) error {
			user := string(key)
			msg := string(value)
			appCtx.Log.Infow("Processing message", "user", user, "msg", msg)
			return nil
		})

		if runErr != nil {
			appCtx.Log.Errorw("Consumer stopped", "error", runErr)
		}
	}()

	//------------------------------------------------------------
	// Ловим сигнал остановки
	//------------------------------------------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	appCtx.Log.Info("🛑 Shutdown signal received")

	//------------------------------------------------------------
	// Отправляем cancel() всем горутинам
	//------------------------------------------------------------
	cancel()

	//------------------------------------------------------------
	// Завершаем Kafka
	//------------------------------------------------------------
	if err := consumer.Close(); err != nil {
		appCtx.Log.Warnw("Kafka consumer close error", "error", err)
	} else {
		appCtx.Log.Info("Kafka consumer closed")
	}

	//------------------------------------------------------------
	// Graceful stop gRPC
	//------------------------------------------------------------
	if grpcServer != nil {
		// GracefulStop не принимает контекст; оборачиваем в горутину, чтобы не блокировать поток
		done := make(chan struct{})
		go func() {
			appCtx.Log.Info("gRPC: GracefulStop called")
			grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			appCtx.Log.Info("gRPC ping stopped gracefully")
		case <-time.After(10 * time.Second):
			appCtx.Log.Warn("gRPC graceful stop timeout, forcing Stop()")
			grpcServer.Stop()
		}
	}

	//------------------------------------------------------------
	// Дожидаемся завершения всех горутин
	//------------------------------------------------------------
	wg.Wait()
	appCtx.Log.Infow("✅ All servers stopped. Bye!")
}

// initDB подключает базу данных и возвращает gorm.DB
func initDB(appCtx *ctx.Context) *gorm.DB {
	dsn := appCtx.Cfg.PostgresDsn
	database, err := db.Connect(dsn)
	if err != nil {
		appCtx.Log.Fatalw("Failed to connect database", "error", err)
	}

	appCtx.Log.Infow("Connected to database", "dsn", dsn)

	return database
}

// startHTTPServer запускает Gin и корректно останавливает его при ctx.Done()
func startHTTPServer(ctx context.Context, appCtx *ctx.Context, wg *sync.WaitGroup, port string) {
	defer wg.Done()

	router := http.SetupRouter(appCtx)
	srv := &netHttp.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	appCtx.Log.Infow("HTTP Server starting", code.Port, port)

	// Запускаем сервер в отдельной горутине для graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, netHttp.ErrServerClosed) {
			appCtx.Log.Errorw("HTTP server error", code.Error, err)
		}
	}()

	// Ожидаем отмены контекста
	<-ctx.Done()
	appCtx.Log.Info("HTTP Server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appCtx.Log.Errorw("HTTP graceful shutdown failed", code.Error, err)
	} else {
		appCtx.Log.Info("HTTP Server stopped gracefully")
	}
}
