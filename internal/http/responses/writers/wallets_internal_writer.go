package writers

import (
	"context"
	"encoding/json"
	"net/http"

	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http/responses"
	"payment-engine/pkg/errors"

	"github.com/recodextech/container"
	"github.com/recodextech/krouter"
)

type WalletsInternalWriter struct {
	log adaptors.Logger
}

func (w *WalletsInternalWriter) Response(_ context.Context, rw http.ResponseWriter, _ *http.Request,
	payload krouter.HttpPayload,
) error {
	var err error
	out := payload.Body.(responses.CreateInternalWalletResponse)

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	resBody, err := json.Marshal(out)
	if err != nil {
		return ResponseWriterError{errors.Wrap(err, errorWritingResponse)}
	}
	_, err = rw.Write(resBody)
	if err != nil {
		return ResponseWriterError{errors.Wrap(err, errorWritingResponse)}
	}

	return nil
}

func (w *WalletsInternalWriter) Init(container container.Container) error {
	w.log = container.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(
		`responses.internal_wallet.create`))
	return nil
}
