package request

// SuccessPaymentRequest represents the request to mark a payment as successful
type SuccessPaymentRequest struct {
	JobID     string `json:"-"`
	PaymentID string `json:"-"`
	AccountID string `json:"-"`
}
