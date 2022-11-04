package logging

type Env string

const (
	Production  Env = "prod"
	Development Env = "dev"
)

var (
	defaultLoggingLevel Level
	defaultLogger       Logger
)

func init() {
	defaultLogger = GetDevelopmentLogger()
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

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	With(args ...interface{}) Logger

	// Sync flushes any buffered log entries.
	Sync() error

	// SetLevel mean set default logger Level
	SetLevel(level Level)
}
