package dao

type Recommendation struct {
	OriginalId string `json:"original_id"`
	InternalId string `json:"internal_id"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	NetworkId  string `json:"network_id"`
	Address    string `json:"address"`
}

type Recommendations []Recommendation
