package logz

import (
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerx struct {
	*zap.Logger
}

func (l loggerx) Printf(msg string, args ...interface{}) {
	if l.Logger.Level() >= zap.InfoLevel {
		l.Logger.Info(fmt.Sprintf(msg, args...))
	}
}

func (l loggerx) Debugf(msg string, args ...interface{}) {
	if l.IsDebug() {
		l.Logger.Debug(fmt.Sprintf(msg, args...))
	}
}

func (l loggerx) IsDebug() bool {
	return l.Logger.Level() >= zap.DebugLevel
}

type nopLogger struct {
	*zap.Logger
}

func (l nopLogger) Debugf(msg string, args ...interface{}) {
}

func (l nopLogger) IsDebug() bool {
	return false
}

var Logger *zap.Logger
var LoggerX loggerx
var NopLogger nopLogger

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
	NopLogger.Logger = zap.NewNop()

	return nil
}
