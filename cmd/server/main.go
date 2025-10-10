package main

import (
	"vado_server/internal/db"
	"vado_server/internal/handlers"
	"vado_server/internal/util"

	"time"
	"vado_server/internal/appcontext"
	"vado_server/internal/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load(".env")
	appCtx := appcontext.NewAppContext(initLogger())
	appCtx.Log.Infow("Start vado-server.", "time", time.Now().Format("2006-01-02 15:04:05"))
	database := initDB(appCtx)
	appCtx.DB = database
	//defer func(database *gorm.DB) {
	//	_ = database.Close()
	//}(database)

	startServer(appCtx, util.GetEnv("PORT"))
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

func startServer(cxt *appcontext.AppContext, port string) {
	gin.SetMode(util.GetEnv("GIN_MODE"))
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	_ = r.SetTrustedProxies(nil)
	r.Static("/static", "./internal/static")
	r.LoadHTMLGlob("internal/templates/*")

	r.GET("/", handlers.ShowIndex)

	r.GET("/tasks", handlers.ShowTasks(cxt))
	r.POST("/tasks", handlers.AddTask(cxt))

	r.GET("/admin", handlers.ShowAdmin(cxt))
	r.POST("/admin", handlers.AddUser(cxt))

	//r.GET("/tasks", handlers.GetTasks(cxt))
	cxt.Log.Infow("Server starting", "port", port)
	if err := r.Run(":" + port); err != nil {
		cxt.Log.Fatalw("Server failed", "error", err)
	}
}
