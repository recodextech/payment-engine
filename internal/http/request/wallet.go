package request

import (
	"encoding/json"
	"payment-engine/pkg/errors"
)

// CreateWalletRequest represents the wallet creation request
type CreateWalletRequest struct {
	AccountID  string `json:"-"`
	WalletType string `json:"wallet_type" validate:"required,oneof=CASH POINTS CREDIT"`
}

func (c CreateWalletRequest) Encode(data interface{}) ([]byte, error) {
	return nil, nil
}

func (c CreateWalletRequest) Decode(data []byte) (interface{}, error) {
	req := CreateWalletRequest{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode create wallet request")
	}
	return req, nil
}
