package request

import "encoding/json"

// TransactionEntry represents a single transaction entry
type TransactionEntry struct {
	Amount  float64 `json:"amount" validate:"required,gt=0"`
	PayeeID string  `json:"payee_id" validate:"required"`
	// WalletType restricts which wallet type the credit is allowed to land in.
	// Accepted values match domain.WalletType: CASH, POINTS, CARD.
	WalletType string `json:"wallet_type" validate:"required,oneof=CASH POINTS CARD"`
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
	// ReasonCode is mandatory and ties every transfer to an auditable business
	// event. See domain.ReasonCode for accepted values.
	ReasonCode string `json:"reason_code" validate:"required,oneof=JOB_COMPLETED FEE_DEDUCTION REFUND POINTS_REDEMPTION"`
	AccountID  string `json:"-"`
	JobID      string `json:"-"`
}

func (c CreatePaymentRequest) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (c CreatePaymentRequest) Decode(data []byte) (interface{}, error) {
	req := CreatePaymentRequest{}
	err := json.Unmarshal(data, &req)
	return req, err
}
