package logz

import (
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerExtension struct {
	*zap.Logger
}

func (l LoggerExtension) Printf(msg string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, args...))
}

var Logger *zap.Logger
var LoggerX LoggerExtension

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
	LoggerX.Logger = Logger
}
