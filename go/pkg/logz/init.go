package logz

import (
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	zConfig := zap.NewProductionConfig()
	zConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zConfigLevel, err := zap.ParseAtomicLevel(config.Values.LogLevel)
	if err != nil {
		panic(err)
	}
	zConfig.Level = zConfigLevel

	sLogger, err := zConfig.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	)
	if err != nil {
		panic(err)
	}

	Logger = sLogger
}
