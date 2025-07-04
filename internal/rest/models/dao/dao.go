package dao

import (
	"time"

	"github.com/google/uuid"
)

type Categories []string

type Strategy struct {
	Name    string                 `json:"name"`
	Network string                 `json:"network"`
	Params  map[string]interface{} `json:"params"`
}

type Strategies []Strategy

type Treasury struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Network string `json:"network"`
}

type Treasuries []Treasury

type Voting struct {
	Delay       uint64  `json:"delay"`
	Period      uint64  `json:"period"`
	Type        string  `json:"type"`
	Quorum      float32 `json:"quorum"`
	Blind       bool    `json:"blind"`
	HideAbstain bool    `json:"hide_abstain"`
	Privacy     string  `json:"privacy"`
	Aliased     bool    `json:"aliased"`
}

type Dao struct {
	ID                 uuid.UUID  `json:"id"`
	Alias              string     `json:"alias"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Name               string     `json:"name"`
	Private            bool       `json:"private"`
	About              string     `json:"about"`
	Avatar             string     `json:"avatar"`
	Terms              string     `json:"terms"`
	Location           string     `json:"location"`
	Website            string     `json:"website"`
	Twitter            string     `json:"twitter"`
	Github             string     `json:"github"`
	Coingecko          string     `json:"coingecko"`
	Email              string     `json:"email"`
	Network            string     `json:"network"`
	Symbol             string     `json:"symbol"`
	Skin               string     `json:"skin"`
	Domain             string     `json:"domain"`
	Strategies         Strategies `json:"strategies"`
	Voting             Voting     `json:"voting"`
	Categories         Categories `json:"categories"`
	Treasures          Treasuries `json:"treasures"`
	FollowersCount     uint64     `json:"followers_count"`
	ProposalsCount     uint64     `json:"proposals_count"`
	Guidelines         string     `json:"guidelines"`
	Template           string     `json:"template"`
	ParentID           string     `json:"parent_id"`
	ActivitySince      uint64     `json:"activity_since"`
	VotersCount        uint64     `json:"voters_count"`
	ActiveVotes        uint64     `json:"active_votes"`
	ActiveProposalsIDs []string   `json:"active_proposals_ids"`
	Verified           bool       `json:"verified"`
	PopularityIndex    float64    `json:"popularity_index"`
	TokenExist         bool       `json:"token_exist"`
	TokenSymbol        string     `json:"token_symbol"`
	FungibleID         string     `json:"fungible_id"`
}
