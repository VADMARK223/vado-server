package appcontext

import (
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AppContext struct {
	Log *zap.SugaredLogger
	DB  *sql.DB
}

func NewAppContext(log *zap.SugaredLogger) *AppContext {
	return &AppContext{Log: log}
}
