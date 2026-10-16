package models

import "time"

type Wallet struct {
	ID        string  `json:"id"`
	AccountID string  `json:"-"`
	Type      string  `json:"type"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"status"`
}

// WalletTotals sums a wallet's ledger: what was credited to it (earned)
// and what was debited from it (paid). Held is what in-progress payments
// will debit from it once their jobs complete.
type WalletTotals struct {
	Credited float64
	Debited  float64
	Held     float64
}

// WalletTransaction is one ledger entry on a wallet. Amount is positive for a
// credit and negative for a debit.
type WalletTransaction struct {
	ID           string
	PaymentID    string
	JobID        string
	WalletType   string
	Amount       float64
	BalanceAfter float64
	CreatedAt    time.Time
}
