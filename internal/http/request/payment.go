package request

import "encoding/json"

// TransactionEntry represents a single credit to a payee wallet
type TransactionEntry struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	// PayeeID is the wallet to credit. When empty, the wallet of PayeeAccountID
	// with the same type as the payer wallet is credited.
	PayeeID        string `json:"payee_id"`
	PayeeAccountID string `json:"payee_account_id"`
	// WalletType restricts which wallet type the credit is allowed to land in.
	// Accepted values match domain.WalletType: CASH, POINTS, CARD.
	WalletType string `json:"wallet_type" validate:"omitempty,oneof=CASH POINTS CARD"`
}

// CreatePaymentRequest represents the payment creation request
type CreatePaymentRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Type   string  `json:"type" validate:"required,oneof=PAYMENT REFUND"`
	// PayerID is the wallet debited by Amount when the payment succeeds.
	PayerID            string             `json:"payer_id" validate:"required"`
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
