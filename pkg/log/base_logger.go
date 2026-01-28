package log

import (
	"context"
	"payment-engine/internal/app"
	"payment-engine/internal/domain/application"

	"github.com/google/uuid"
	"github.com/recodextech/container"
	"github.com/tryfix/log"
	traceableCtx "github.com/tryfix/traceable-context"
)

type BaseLogger struct {
	log.Logger
}

func (l *BaseLogger) Init(c container.Container) error {
	con := c.GetGlobalConfig(application.ModuleBaseLogger).(*LoggerConf)

	logger := log.NewLog(log.WithOutput(log.OutJson))

	// when in debug mode use the text logger
	if app.DebugMode() {
		logger = log.NewLog(log.WithOutput(log.OutText), log.WithColors(con.Colors))
	}

	l.Logger = logger.Log(
		log.WithLevel(log.Level(con.Level)),
		log.Prefixed(`application`),
		log.WithSkipFrameCount(3), // nolint
		log.WithFilePath(con.FilePath),
		log.WithCtxTraceExtractor(func(ctx context.Context) string {
			if trace := traceableCtx.FromContext(ctx); trace != uuid.Nil {
				return trace.String()
			}
			return ""
		}),
	)

	return nil
}
