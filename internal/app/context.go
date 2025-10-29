package app

import (
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Context struct {
	Log *zap.SugaredLogger
	DB  *gorm.DB
}

func NewAppContext(log *zap.SugaredLogger) *Context {
	return &Context{Log: log}
}
