package handlers

import (
	"errors"
	"net/http"
	"payment-engine/internal/http/responses"
	"payment-engine/internal/services"
)

// format formats the given error in the output format.
func format(err error) responses.ErrorResponse {
	switch errType := err.(type) {
	case services.ServiceError:
		status := errType.Status
		if status == 0 {
			status = http.StatusBadRequest
		}
		return responses.ErrorResponse{
			Description: errType.Msg,
			Code:        errType.Code,
			Status:      status,
			Details:     errType.Details,
		}
	default:
		// check whether the embedded error is of a known error type
		e := errors.Unwrap(err)
		if e == nil {
			return responses.ErrorResponse{
				Description: "unknown error",
				Status:      http.StatusInternalServerError,
			}
		}

		return format(e)
	}
}
