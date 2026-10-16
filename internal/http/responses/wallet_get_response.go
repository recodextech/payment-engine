package responses

import (
	"payment-engine/internal/domain/models"

	"github.com/recodextech/api-definitions/events"
)

type GetWalletsResponse struct {
	Wallets []WalletPayload `json:"wallets"`
}

type WalletPayload struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
	Status  string  `json:"status"`
	// TotalCredited is everything paid into the wallet (worker earnings)
	TotalCredited float64 `json:"totalCredited"`
	// TotalDebited is everything paid out of the wallet (contractor payments)
	TotalDebited float64 `json:"totalDebited"`
	// HeldAmount is reserved by in-progress payments until their jobs end
	HeldAmount float64 `json:"heldAmount"`
	// AvailableBalance is the balance less the held amount; it may be negative
	AvailableBalance   float64 `json:"availableBalance"`
	PaymentInformation struct {
		Method    string            `json:"method"`
		Reference map[string]string `json:"reference"`
	} `json:"walletInformation"`
}

// NewGetWalletsResponse builds the wallets list; totals are keyed by wallet ID
func NewGetWalletsResponse(wallets []events.AccountWalletEvent, totals map[string]models.WalletTotals) GetWalletsResponse {
	payload := make([]WalletPayload, 0, len(wallets))
	for _, w := range wallets {
		wallet := WalletPayload{
			ID:            w.Payload.ID,
			Type:          w.Payload.Type,
			Balance:       w.Payload.Balance,
			Status:        w.Payload.Status,
			TotalCredited: totals[w.Payload.ID].Credited,
			TotalDebited:  totals[w.Payload.ID].Debited,

			HeldAmount:       totals[w.Payload.ID].Held,
			AvailableBalance: w.Payload.Balance - totals[w.Payload.ID].Held,
		}
		wallet.PaymentInformation.Method = w.Payload.PaymentInfromation.Method
		wallet.PaymentInformation.Reference = w.Payload.PaymentInfromation.Reference
		payload = append(payload, wallet)
	}
	return GetWalletsResponse{
		Wallets: payload,
	}
}
