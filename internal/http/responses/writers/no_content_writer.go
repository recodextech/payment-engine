package writers

import (
	"context"
	"net/http"

	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/application"

	"github.com/recodextech/container"
	"github.com/recodextech/krouter"
)

type NoContentWriter struct {
	log adaptors.Logger
}

func (w *NoContentWriter) Response(_ context.Context, rw http.ResponseWriter, _ *http.Request,
	_ krouter.HttpPayload,
) error {
	rw.WriteHeader(http.StatusNoContent)
	return nil
}

func (w *NoContentWriter) Init(c container.Container) error {
	w.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(
		`responses.no-content`))
	return nil
}
