package logz

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	zConfig := zap.NewProductionConfig()
	zConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	sLogger, err := zConfig.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	)
	if err != nil {
		panic(err)
	}

	Logger = sLogger
}
