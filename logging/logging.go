package logging

import (
	"fmt"

	"github.com/spf13/viper"
)

type Env string

const (
	Production  Env = "prod"
	Development Env = "dev"
)

var (
	defaultLogger       Logger
	developmentLogger   Logger
	productionLogger    Logger
	defaultLoggingLevel Level
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
}

func SetLogger(l Logger) {
	defaultLogger = l
}

func GetDefaultLogger() Logger {
	return defaultLogger
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

func Sync() error {
	return defaultLogger.Sync()
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	With(args ...interface{}) Logger

	// Sync flushes any buffered log entries.
	Sync() error
}
