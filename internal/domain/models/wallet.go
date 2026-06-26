package models

type Wallet struct {
	ID        string  `json:"id"`
	AccountID string  `json:"-"`
	Type      string  `json:"type"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"status"`
}
