package validator

type ValidatorError struct {
	error
	Message string
}

const (
	ErrMsgInvalidAccountID = "invalid account ID"
)
