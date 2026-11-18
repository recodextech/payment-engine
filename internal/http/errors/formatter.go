package errors

import (
	"errors"
	"net/http"
	response "payment-engine/internal/http/responses"
	"payment-engine/internal/services"
)

// format formats the given error in the output format.
func format(err error) response.ErrorResponse {
	switch errType := err.(type) {
	case services.ServiceError:
		status := errType.Status
		if status == 0 {
			status = http.StatusBadRequest
		}
		return response.ErrorResponse{
			Description: errType.Msg,
			Code:        errType.Code,
			Status:      status,
			Details:     errType.Details,
		}
	default:
		// check whether the embedded error is of a known error type
		e := errors.Unwrap(err)
		if e == nil {
			return response.ErrorResponse{
				Description: "unknown error",
				Status:      http.StatusInternalServerError,
			}
		}

		return format(e)
	}
}
