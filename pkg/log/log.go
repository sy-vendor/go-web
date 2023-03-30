package log

import (
	"os"
	"sync"

	"github.com/google/wire"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ProviderSet = wire.NewSet(NewLogger)
var once sync.Once
var logger *zap.Logger

func NewLogger() *zap.Logger {
	once.Do(func() {
		ecsEncoder := ecszap.NewDefaultEncoderConfig()
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    500, //megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		})
		writer := zapcore.AddSync(w)
		defaultLogLevel := zapcore.DebugLevel
		core := zapcore.NewTee(
			ecszap.NewCore(ecsEncoder, writer, defaultLogLevel),
			ecszap.NewCore(ecsEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		logger = logger.With(zap.String("app", "web")).
			With(zap.String("environment", "test"))
		logger.Info("logger initialized")
	})

	return logger
}

func Close() {
	_ = logger.Sync()
}
