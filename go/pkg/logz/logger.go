package logz

import (
	"fmt"

	"github.com/tartale/go/pkg/jsontime"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerExtension struct {
	*zap.Logger
}

func (l LoggerExtension) Printf(msg string, args ...interface{}) {
	if l.Logger.Level() >= zap.InfoLevel {
		l.Logger.Info(fmt.Sprintf(msg, args...))
	}
}

func (l LoggerExtension) Debugf(msg string, args ...interface{}) {
	if l.IsDebug() {
		l.Logger.Debug(fmt.Sprintf(msg, args...))
	}
}

func (l LoggerExtension) IsDebug() bool {
	return l.Logger.Level() >= zap.DebugLevel
}

var Logger *zap.Logger
var LoggerX LoggerExtension

func InitLoggers() error {

	zConfig := zap.NewProductionConfig()
	zConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zConfigLevel, err := zap.ParseAtomicLevel(config.Values.LogLevel)
	if err != nil {
		return err
	}
	zConfig.Level = zConfigLevel

	sLogger, err := zConfig.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	)
	if err != nil {
		return err
	}

	Logger = sLogger
	LoggerX.Logger = Logger
	jsontime.Logger = LoggerX

	return nil
}
