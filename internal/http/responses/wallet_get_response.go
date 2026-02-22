package responses

import "github.com/recodextech/api-definitions/events"

type GetWalletsResponse struct {
	Wallets []WalletPayload `json:"wallets"`
}

type WalletPayload struct {
	ID        string  `json:"id"`
	AccountID string  `json:"account_id"`
	Type      string  `json:"type"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"status"`
}

func NewGetWalletsResponse(wallets []events.AccountWalletEvent) GetWalletsResponse {
	var payload []WalletPayload
	for _, w := range wallets {
		payload = append(payload, WalletPayload{
			ID:        w.Payload.ID,
			AccountID: w.Payload.AccountID,
			Type:      w.Payload.Type,
			Balance:   w.Payload.Balance,
			Status:    w.Payload.Status,
		})
	}
	return GetWalletsResponse{
		Wallets: payload,
	}
}
