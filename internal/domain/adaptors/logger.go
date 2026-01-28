package adaptors

import "context"

type LogLevel string

const (
	FATAL LogLevel = `FATAL`
	ERROR LogLevel = `ERROR`
	WARN  LogLevel = `WARN`
	INFO  LogLevel = `INFO`
	DEBUG LogLevel = `DEBUG`
	TRACE LogLevel = `TRACE`
)

type LoggerOptions map[string]interface{}

func NewLoggerOptions() LoggerOptions {
	return map[string]interface{}{}
}

type LoggerOption func(LoggerOptions)

/*
Add a prefix to the existing log instance
*/
func LoggerPrefixed(prefix string) LoggerOption {
	return func(options LoggerOptions) {
		options[`prefix`] = prefix
	}
}

/*
Set the log level
*/
func LoggerWithLevel(level LogLevel) LoggerOption {
	return func(options LoggerOptions) {
		options[`level`] = level
	}
}

/*
implement Logger interface for logging within the application
*/
type Logger interface {
	Fatal(message interface{}, params ...interface{})
	Error(message interface{}, params ...interface{})
	Warn(message interface{}, params ...interface{})
	Debug(message interface{}, params ...interface{})
	Info(message interface{}, params ...interface{})
	Trace(message interface{}, params ...interface{})
	FatalContext(ctx context.Context, message interface{}, params ...interface{})
	ErrorContext(ctx context.Context, message interface{}, params ...interface{})
	WarnContext(ctx context.Context, message interface{}, params ...interface{})
	DebugContext(ctx context.Context, message interface{}, params ...interface{})
	InfoContext(ctx context.Context, message interface{}, params ...interface{})
	TraceContext(ctx context.Context, message interface{}, params ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	NewLog(...LoggerOption) Logger
	Params(key, value string) func(k, v string) string
}
