package responses

import (
	"time"

	"payment-engine/internal/domain/models"
)

const (
	TransactionDirectionCredit = "CREDIT"
	TransactionDirectionDebit  = "DEBIT"
)

type WalletTransactionsResponse struct {
	WalletID     string                      `json:"walletId"`
	Transactions []WalletTransactionResponse `json:"transactions"`
}

type WalletTransactionResponse struct {
	ID        string `json:"id"`
	PaymentID string `json:"paymentId"`
	JobID     string `json:"jobId,omitempty"`
	Type      string `json:"type"`
	// Direction is CREDIT (money in) or DEBIT (money out); Amount is always positive
	Direction    string  `json:"direction"`
	Amount       float64 `json:"amount"`
	BalanceAfter float64 `json:"balanceAfter"`
	CreatedAt    string  `json:"createdAt"`
}

func NewWalletTransactionsResponse(walletID string, transactions []models.WalletTransaction) WalletTransactionsResponse {
	res := WalletTransactionsResponse{
		WalletID:     walletID,
		Transactions: make([]WalletTransactionResponse, 0, len(transactions)),
	}
	for _, t := range transactions {
		direction, amount := TransactionDirectionCredit, t.Amount
		if t.Amount < 0 {
			direction, amount = TransactionDirectionDebit, -t.Amount
		}
		res.Transactions = append(res.Transactions, WalletTransactionResponse{
			ID:           t.ID,
			PaymentID:    t.PaymentID,
			JobID:        t.JobID,
			Type:         t.WalletType,
			Direction:    direction,
			Amount:       amount,
			BalanceAfter: t.BalanceAfter,
			CreatedAt:    t.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	return res
}
