package dao

import (
	"time"
)

type DelegatesResponse struct {
	Delegates []Delegate `json:"delegates"`
	Total     int32      `json:"total"`
}

type Delegate struct {
	Address               string  `json:"address"`
	ENSName               string  `json:"ens_name"`
	DelegatorCount        int32   `json:"delegator_count"`
	PercentOfDelegators   float64 `json:"percent_of_delegators"`
	VotingPower           float64 `json:"voting_power"`
	PercentOfVotingPower  float64 `json:"percent_of_voting_power"`
	About                 string  `json:"about"`
	Statement             string  `json:"statement"`
	VotesCount            int32   `json:"votes_count"`
	CreatedProposalsCount int32   `json:"created_proposals_count"`
}

type DelegateProfile struct {
	Address              string                `json:"address"`
	VotingPower          float64               `json:"voting_power"`
	IncomingPower        float64               `json:"incoming_power"`
	OutgoingPower        float64               `json:"outgoing_power"`
	PercentOfVotingPower float64               `json:"percent_of_voting_power"`
	PercentOfDelegators  float64               `json:"percent_of_delegators"`
	Delegates            []ProfileDelegateItem `json:"delegates"`
	Expiration           *time.Time            `json:"expiration,omitempty"`
}

type ProfileDelegateItem struct {
	Address        string  `json:"address"`
	ENSName        string  `json:"ens_name"`
	Weight         float64 `json:"weight"`
	DelegatedPower float64 `json:"delegated_power"`
}
