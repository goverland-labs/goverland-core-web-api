package dao

import "time"

type TokenInfo struct {
	Name                  string           `json:"name"`
	Symbol                string           `json:"symbol"`
	TotalSupply           float64          `json:"total_supply"`
	CirculatingSupply     float64          `json:"circulating_supply"`
	MarketCap             float64          `json:"market_cap"`
	FullyDilutedValuation float64          `json:"fully_diluted_valuation"`
	Price                 float64          `json:"price"`
	FungibleID            string           `json:"fungible_id"`
	Chains                []TokenChainInfo `json:"chains"`
}

type TokenChainInfo struct {
	ChainID  string `json:"chain_id"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	IconURL  string `json:"icon_url"`
	Address  string `json:"address"`
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
