package handlers

import (
	"errors"
	"payment-engine/internal/http/responses"
	"payment-engine/internal/services"
	"net/http"
)

// format formats the given error in the output format.
func format(err error) responses.ErrorResponse {
	switch errType := err.(type) {
	case services.ServiceError:
		return responses.ErrorResponse{
			Description: errType.Msg,
			Code:        errType.Code,
			Status:      http.StatusBadRequest,
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
