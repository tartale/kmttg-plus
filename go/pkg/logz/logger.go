package logz

import (
	"fmt"

	gologz "github.com/tartale/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerx struct {
	*zap.SugaredLogger
}

func (l loggerx) Printf(msg string, args ...interface{}) {
	if l.SugaredLogger.Level() >= zap.InfoLevel {
		l.SugaredLogger.Info(fmt.Sprintf(msg, args...))
	}
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
	LoggerX.SugaredLogger = Logger.Sugar()
	NopLogger.Logger = zap.NewNop()

	return nil
}

func InitThirdPartyLoggers() error {

	gologz.SetLoggerForName("github.com/tartale/go/pkg/jsontime", NopLogger)
	gologz.SetLoggerForName("github.com/tartale/go/pkg/generics", LoggerX)

	return nil
}
