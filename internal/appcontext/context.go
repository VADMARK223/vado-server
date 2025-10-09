package appcontext

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AppContext struct {
	Log *zap.SugaredLogger
	DB  *sql.DB
}

func NewAppContext(log *zap.SugaredLogger) *AppContext {
	app := &AppContext{Log: log}
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalw("failed to connect to database", "error", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalw("cannot ping database", "error", err)
	}

	app.DB = db
	log.Infow("Database connected")

	return app
}
