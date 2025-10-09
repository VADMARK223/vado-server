package appcontext

import (
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Log *zap.SugaredLogger
	DB  *gorm.DB
}

func NewAppContext(log *zap.SugaredLogger) *AppContext {
	return &AppContext{Log: log}
}
