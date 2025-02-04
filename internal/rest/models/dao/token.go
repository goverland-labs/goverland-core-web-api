package dao

import "time"

type TokenInfo struct {
	Name                  string  `json:"name"`
	Symbol                string  `json:"symbol"`
	TotalSupply           float64 `json:"total_supply"`
	CirculatingSupply     float64 `json:"circulating_supply"`
	MarketCap             float64 `json:"market_cap"`
	FullyDilutedValuation float64 `json:"fully_diluted_valuation"`
	Price                 float64 `json:"price"`
}

type TokenChart struct {
	Price        float64 `json:"price"`
	PriceChanges float64 `json:"price_changes"`
	Points       []Point `json:"points"`
}

type Point struct {
	Time  time.Time `json:"time"`
	Price float64   `json:"price"`
}
