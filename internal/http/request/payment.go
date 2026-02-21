package request

import "encoding/json"

// TransactionEntry represents a single transaction entry
type TransactionEntry struct {
	Amount  float64 `json:"amount" validate:"required,gt=0"`
	PayeeID string  `json:"payee_id" validate:"required"`
}

// CreatePaymentRequest represents the payment creation request
type CreatePaymentRequest struct {
	Amount             float64            `json:"amount" validate:"required,gt=0"`
	Type               string             `json:"type" validate:"required,oneof=PAYMENT REFUND"`
	TransactionEntries []TransactionEntry `json:"transaction_entries" validate:"required,min=1"`
	PaymentInfromation struct {
		Method    string            `json:"method"`
		Reference map[string]string `json:"reference"`
	} `json:"payment_information"`
	AccountID string `json:"-"`
	JobID     string `json:"-"`
}

func (c CreatePaymentRequest) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (c CreatePaymentRequest) Decode(data []byte) (interface{}, error) {
	req := CreatePaymentRequest{}
	err := json.Unmarshal(data, &req)
	return req, err
}
