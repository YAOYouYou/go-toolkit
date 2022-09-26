package logging

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	flushLogs           func() error
	defaultLogger       Logger
	developmentLogger   Logger
	productionLogger    Logger
	defaultLoggingLevel Level
)

type Level = zapcore.Level

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

type Env string

const (
	Production  Env = "prod"
	Development Env = "dev"
)

func init() {
	lvl := viper.GetInt("logging-level")
	defaultLoggingLevel = Level(lvl)
	env := Env(viper.GetString("env"))
	switch env {
	case Production:
		defaultLogger = GetProductionLogger()
	case Development:
		defaultLogger = GetDevelopmentLogger()
	default:
		fmt.Printf("unknow env: {%s}, init Development Logger for project.", env)
		defaultLogger = GetDevelopmentLogger()
	}
	sugaredLogger := defaultLogger.(*zap.SugaredLogger)
	flushLogs = sugaredLogger.Sync
}

func GetProductionLogger() Logger {
	if productionLogger != nil {
		return productionLogger
	}
	encoder := getEncoder()
	writeSyncer := getWriteSyncer()
	levelEnabler := zap.LevelEnablerFunc(func(level Level) bool {
		return level >= defaultLoggingLevel
	})
	core := zapcore.NewCore(encoder, writeSyncer, levelEnabler)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	productionLogger = zapLogger.Sugar()
	return productionLogger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := getProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "datetime",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getWriteSyncer() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stderr)
}

func GetDevelopmentLogger() Logger {
	if developmentLogger != nil {
		return developmentLogger
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(defaultLoggingLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLogger, _ := cfg.Build()
	zapLogger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
	developmentLogger = zapLogger.Sugar()
	return developmentLogger
}

func GetDefaultProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{}
}

func GetDefaultLogger() Logger {
	return defaultLogger
}

func LogLevel() string {
	return defaultLoggingLevel.String()
}

func CreateLoggerAsLocalFile(localFilePath string, logLevel Level) (logger Logger, flush func() error, err error) {
	if len(localFilePath) == 0 {
		return nil, nil, errors.New("invalid local logger path")
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   localFilePath,
		MaxSize:    100, // megabytes
		MaxBackups: 2,
		MaxAge:     15, // days
	}
	encoder := getEncoder()
	ws := zapcore.AddSync(lumberJackLogger)
	zapcore.Lock(ws)
	levelEnabler := zap.LevelEnablerFunc(func(level Level) bool {
		return level >= logLevel
	})
	core := zapcore.NewCore(encoder, ws, levelEnabler)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger = zapLogger.Sugar()
	flush = zapLogger.Sync
	return
}

func Cleanup() {
	if flushLogs != nil {
		_ = flushLogs()
	}
}

func Error(err error) {
	if err != nil {
		defaultLogger.Errorf("error occurs during runtime, %v", err)
	}
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	With(args ...interface{}) *zap.SugaredLogger
}
