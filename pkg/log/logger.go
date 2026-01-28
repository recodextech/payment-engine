package log

import (
	"context"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http/request"

	"github.com/recodextech/api-definitions/pkg/uuid"
	"github.com/recodextech/container"
	"github.com/tryfix/log"
)

type Logger struct {
	logger log.Logger
}

func (l *Logger) Init(container container.Container) error {
	l.logger = container.Resolve(application.ModuleBaseLogger).(log.Logger).NewLog(
		log.WithCtxExtractor(func(ctx context.Context) []interface{} {
			var params []interface{}
			if ctx.Value(request.TraceID) != nil {
				params = append(params, string(request.TraceID)+":"+ctx.Value(request.TraceID).(uuid.UUID).String())
			}
			if ctx.Value(request.AccountID) != nil {
				params = append(params, string(request.AccountID)+":"+ctx.Value(request.AccountID).(uuid.UUID).String())
			}
			if ctx.Value(request.UserID) != nil {
				params = append(params, string(request.UserID)+":"+ctx.Value(request.UserID).(uuid.UUID).String())
			}

			return params
		})).NewLog(log.WithCtxTraceExtractor(func(ctx context.Context) string {
		if ctx.Value(request.TraceID) != nil {
			return ctx.Value(request.TraceID).(uuid.UUID).String()
		}

		return uuid.New().String()
	}))

	return nil
}

func (l *Logger) Params(key, value string) func(k, v string) string {
	return func(k, v string) string {
		return key + ":" + value
	}
}

func newLogger(log log.Logger) adaptors.Logger {
	return &Logger{logger: log}
}

func (l *Logger) Fatal(message interface{}, params ...interface{}) {
	l.logger.Fatal(message, params...)
}

func (l *Logger) Error(message interface{}, params ...interface{}) {
	l.logger.Error(message, params...)
}

func (l *Logger) Warn(message interface{}, params ...interface{}) {
	l.logger.Warn(message, params...)
}

func (l *Logger) Debug(message interface{}, params ...interface{}) {
	l.logger.Debug(message, params...)
}

func (l *Logger) Info(message interface{}, params ...interface{}) {
	l.logger.Info(message, params...)
}

func (l *Logger) Trace(message interface{}, params ...interface{}) {
	l.logger.Trace(message, params...)
}

func (l *Logger) FatalContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.logger.FatalContext(ctx, message, params...)
}

func (l *Logger) ErrorContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.logger.ErrorContext(ctx, message, params...)
}

func (l *Logger) WarnContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.logger.WarnContext(ctx, message, params...)
}

func (l *Logger) DebugContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.logger.DebugContext(ctx, message, params...)
}

func (l *Logger) InfoContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.logger.InfoContext(ctx, message, params...)
}

func (l *Logger) TraceContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.logger.TraceContext(ctx, message, params...)
}

func (l *Logger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *Logger) NewLog(options ...adaptors.LoggerOption) adaptors.Logger {
	optMap := adaptors.NewLoggerOptions()
	for _, opt := range options {
		opt(optMap)
	}

	var opts []log.Option
	for typ, opt := range optMap {
		switch typ {
		case `prefix`:
			opts = append(opts, log.Prefixed(opt.(string)))
		case `level`:
			opts = append(opts, log.WithLevel(log.Level(opt.(adaptors.LogLevel))))
		}
	}

	return newLogger(l.logger.NewLog(opts...))
}
