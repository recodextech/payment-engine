package responses

import "github.com/recodextech/api-definitions/events"

type GetWalletsResponse struct {
	Wallets []WalletPayload `json:"wallets"`
}

type WalletPayload struct {
	ID                 string  `json:"id"`
	Type               string  `json:"type"`
	Balance            float64 `json:"balance"`
	Status             string  `json:"status"`
	PaymentInformation struct {
		Method    string            `json:"method"`
		Reference map[string]string `json:"reference"`
	} `json:"walletInformation"`
}

func NewGetWalletsResponse(wallets []events.AccountWalletEvent) GetWalletsResponse {
	var payload []WalletPayload
	for _, w := range wallets {
		wallet := WalletPayload{
			ID:      w.Payload.ID,
			Type:    w.Payload.Type,
			Balance: w.Payload.Balance,
			Status:  w.Payload.Status,
		}
		wallet.PaymentInformation.Method = w.Payload.PaymentInfromation.Method
		wallet.PaymentInformation.Reference = w.Payload.PaymentInfromation.Reference
		payload = append(payload, wallet)
	}
	return GetWalletsResponse{
		Wallets: payload,
	}
}
