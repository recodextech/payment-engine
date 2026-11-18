package services

import (
	"payment-engine/pkg/errors"
)

// ServiceError represent the type of custom errors returned by the services.
type ServiceError struct {
	error
	Code int
	Msg  string
	// Status is the HTTP status to answer with; 0 means 400 Bad Request.
	Status int
	// Details carries machine readable fields for the client, e.g. remaining job quota.
	Details map[string]any
}

var (
	ErrMsgResNotFound = errors.Msg(errors.CodeResNotFound)
)

// NewServiceError builds a ServiceError answered with the given HTTP status.
func NewServiceError(status, code int, msg string) ServiceError {
	return ServiceError{error: errors.New(msg), Code: code, Msg: msg, Status: status}
}
