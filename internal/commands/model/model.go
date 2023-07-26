package model

type GoogleApiResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type KhugaWalletCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Whitelisted bool `json:"whitelisted"`
	} `json:"data"`
}
