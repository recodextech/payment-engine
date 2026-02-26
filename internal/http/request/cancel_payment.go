package request

// CancelPaymentRequest represents the request to cancel a payment
type CancelPaymentRequest struct {
	JobID     string `json:"-"`
	PaymentID string `json:"-"`
	AccountID string `json:"-"`
}
