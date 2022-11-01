package logging

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	*zap.SugaredLogger
}

func (l *zapLogger) With(args ...interface{}) Logger {
	logger := l.SugaredLogger.With(args...)
	return &zapLogger{SugaredLogger: logger}
}

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
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	productionLogger = &zapLogger{SugaredLogger: logger.Sugar()}
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
	logger, _ := cfg.Build()
	logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
	developmentLogger = &zapLogger{SugaredLogger: logger.Sugar()}
	return developmentLogger
}

func GetDefaultProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{}
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
	l := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger = &zapLogger{SugaredLogger: l.Sugar()}
	flush = logger.(*zapLogger).Sync
	return
}
