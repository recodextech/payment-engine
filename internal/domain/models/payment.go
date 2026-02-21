package models

type Payment struct {
	ID                 string               `json:"id"`
	JobAssignedID      float64              `json:"-"`
	Type               string               `json:"-"` // e.g., "PAYMENT", "REFUND", "REVERSAL"
	Amount             float64              `json:"-"`
	TransactionEntries []TransactionPayload `json:"-"`
	Status             string               `json:"-"`
	PaymentInfromation struct {
		Method    string            `json:"method"`
		Version   int               `json:"version"`
		Reference map[string]string `json:"reference"`
	} `json:"payment_information"`
}

type TransactionPayload struct {
	ID        string  `json:"id"`
	PaymentID string  `json:"-"`
	Type      string  `json:"-"` // e.g., "CREDIT_CARD", "BANK_TRANSFER", "CASH", WALLET", "FEE"
	Amount    float64 `json:"-"`
	Payee     struct {
		ID             string  `json:"id"`
		CurrentBalance float64 `json:"current_balance"`
		SequenceNumber int64   `json:"sequence_number"`
	} `json:"payee"`
	PaymentInfromation struct {
		Method    string            `json:"method"`
		Version   int               `json:"version"`
		Reference map[string]string `json:"reference"`
	} `json:"payment_information"`
}
