package app

import (
	"vado_server/internal/config/config"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Context struct {
	Log *zap.SugaredLogger
	DB  *gorm.DB
	Cfg *config.Config
}

func NewAppContext(log *zap.SugaredLogger) *Context {
	cfg := config.Load()
	return &Context{
		Log: log,
		Cfg: cfg,
	}
}
