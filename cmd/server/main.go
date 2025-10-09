package main

import (
	"net/http"
	"os"
	"time"
	"vado_server/internal/appcontext"
	"vado_server/internal/handler"
	"vado_server/internal/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load(".env")
	appCtx := appcontext.NewAppContext(initLogger())
	initLogger()

	handler.RegisterTaskRoutes(appCtx)

	// POSTGRES_URL is not set

	startServer(appCtx, getPort())
}

func initLogger() *zap.SugaredLogger {
	log, loggerInitErr := logger.Init(true)
	if loggerInitErr != nil {
		panic(loggerInitErr)
	}
	defer func() { _ = log.Sync() }()

	return log
}

func startServer(ctx *appcontext.AppContext, port string) {
	ctx.Log.Infow("Start vado-server.", "time", time.Now().Format("2006-01-02 15:04:05"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello from vado-server"))
	})
	ctx.Log.Infow("Server started", "port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		ctx.Log.Errorw("Server error", "error", err)
		return
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5555"
	}
	return port
}
