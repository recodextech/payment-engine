package responses

type CreateWalletResponse struct {
	ID string `json:"id"`
}

type CreateInternalWalletResponse struct {
	IDs []string `json:"ids"`
}
