package writers

import (
	"context"
	"encoding/json"
	"net/http"

	"payment-engine/pkg/errors"

	"github.com/recodextech/container"
	"github.com/recodextech/krouter"
)

// JSONWriter writes the handler result as a 200 JSON body.
type JSONWriter struct{}

func (w *JSONWriter) Init(_ container.Container) error { return nil }

func (w *JSONWriter) Response(_ context.Context, rw http.ResponseWriter, _ *http.Request, payload krouter.HttpPayload) error {
	body, err := json.Marshal(payload.Body)
	if err != nil {
		return ResponseWriterError{errors.Wrap(err, errorWritingResponse)}
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(body); err != nil {
		return ResponseWriterError{errors.Wrap(err, errorWritingResponse)}
	}
	return nil
}
