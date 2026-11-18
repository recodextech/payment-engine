package errors

import (
	"context"
	"encoding/json"
	"net/http"
	"payment-engine/internal/app"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http/request"
	perrors "payment-engine/pkg/errors"

	"github.com/recodextech/api-definitions/pkg/uuid"
	"github.com/recodextech/container"
)

// Handler contains all error handling, formatting and presenting functionality for the http layer.
type Handler struct {
	log adaptors.Logger
}

// Init initializes a new instance of the error handler.
func (hlr *Handler) Init(c container.Container) error {
	hlr.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(
		adaptors.LoggerPrefixed("ErrorHandler"))

	return nil
}

// Handle handles all errors globally.
func (hlr *Handler) Handle(ctx context.Context, w http.ResponseWriter, e error) error {
	// log the error
	hlr.log.ErrorContext(ctx, e.Error())
	traceID := ctx.Value(request.TraceID).(uuid.UUID)

	resError := format(e)
	resError.Trace = traceID.String()

	if app.DebugMode() {
		resError.Debug = e.Error()
	}

	payload, err := json.Marshal(resError)
	if err != nil {
		return perrors.Wrap(err, "failed to marshal error response")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(resError.Status)

	_, err = w.Write(payload)
	if err != nil {
		return perrors.Wrap(err, "failed to write error response")
	}

	return nil
}
