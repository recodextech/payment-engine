package errors

import (
	"errors"
	response "payment-engine/internal/http/responses"
	"payment-engine/internal/services"
	"net/http"
)

// format formats the given error in the output format.
func format(err error) response.ErrorResponse {
	switch errType := err.(type) {
	case services.ServiceError:
		return response.ErrorResponse{
			Description: errType.Msg,
			Code:        errType.Code,
			Status:      http.StatusBadRequest,
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
