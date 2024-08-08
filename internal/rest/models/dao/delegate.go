package dao

type Delegate struct {
	Address                  string  `json:"address"`
	ENSName                  string  `json:"ens_name"`
	DelegatorCount           int32   `json:"delegator_count"`
	PercentOfDelegators      float64 `json:"percent_of_delegators"`
	VotingPower              float64 `json:"voting_power"`
	PercentOfVotingPower     float64 `json:"percent_of_voting_power"`
	About                    string  `json:"about"`
	Statement                string  `json:"statement"`
	UserDelegatedVotingPower float64 `json:"user_delegated_voting_power"`
	VotesCount               int32   `json:"votes_count"`
	ProposalsCount           int32   `json:"proposals_count"`
	CreateProposalsCount     int32   `json:"create_proposals_count"`
}
