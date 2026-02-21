package models

type Wallet struct {
	ID        string  `json:"id"`
	AccountID string  `json:"-"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"status"`
}
