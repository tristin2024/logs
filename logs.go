//author: jiazujiang
//date: 2023/6/14

package logs

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Std *zap.SugaredLogger

var (
	Env string
	App string = "backend"
)

const (
	envProd = "prod"
	envTest = "test"
)

func init() {
	Env = os.Getenv("Env")
	App = os.Getenv("app")
	var serviceName string
	if (Env == envProd || Env == envTest) && len(App) != 0 {
		serviceName = App
	} else {
		serviceName = "backend"
	}
	Std = newLogging(serviceName)
}

func newLogging(service string) *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var atom zap.AtomicLevel
	var outputPaths []string
	var errorOutputPaths []string

	if Env == envProd || Env == envTest {
		atom = zap.NewAtomicLevelAt(zap.WarnLevel)
		outputPaths = []string{"stdout", "./app.log"}
		errorOutputPaths = []string{"stderr"}
	} else {
		atom = zap.NewAtomicLevelAt(zap.DebugLevel)
		outputPaths = []string{"stdout"}
		errorOutputPaths = []string{"stderr"}
	}

	encoding := os.Getenv("LOG_ENCODING")
	if len(encoding) == 0 {
		encoding = "json"
	}
	config := zap.Config{
		Level:             atom,
		DisableCaller:     false,
		Development:       false,
		DisableStacktrace: true,
		Encoding:          encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  errorOutputPaths,
	}
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	var err error
	log, err := config.Build()
	log = log.Named(service).With(zap.Any("env", Env))
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}
	return log.Sugar()
}
