package appcontext

import "go.uber.org/zap"

type AppContext struct {
	Log *zap.SugaredLogger
}

func NewAppContext(log *zap.SugaredLogger) *AppContext {
	return &AppContext{
		Log: log,
	}
}
