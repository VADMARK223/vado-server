package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(dev bool) *zap.SugaredLogger {
	var encoderCfg zapcore.EncoderConfig
	var encoder zapcore.Encoder
	var minLevel zapcore.Level

	if dev {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // Разными цветами раскрашивает типы ошибок
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		minLevel = zap.DebugLevel // показываем всё
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(encoderCfg)
		minLevel = zap.InfoLevel // скрываем DEBUG
	}

	// Разделяем уровни
	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= minLevel && l < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), errorLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
}
