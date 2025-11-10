package dao

type Delegator struct {
	Address    string `json:"address"`
	ENSName    string `json:"ens_name"`
	TokenValue string `json:"token_value"`
}
