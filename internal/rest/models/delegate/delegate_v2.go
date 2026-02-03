package delegate

import "time"

type (
	TokenValue struct {
		Value    string `json:"value"`
		Symbol   string `json:"symbol"`
		Decimals int32  `json:"decimals"`
	}

	DelegateEntryV2 struct {
		Address               string      `json:"address"`
		EnsName               string      `json:"ens_name"`
		DelegatorCount        *int32      `json:"delegator_count,omitempty"`
		PercentOfDelegators   *float64    `json:"percent_of_delegators,omitempty"`
		PercentOfVotingPower  *float64    `json:"percent_of_voting_power,omitempty"`
		About                 *string     `json:"about,omitempty"`
		Statement             *string     `json:"statement,omitempty"`
		VotesCount            *int32      `json:"votes_count,omitempty"`
		CreatedProposalsCount *int32      `json:"created_proposals_count,omitempty"`
		VotingPower           *float64    `json:"voting_power,omitempty"`
		TokenValue            *TokenValue `json:"token_value,omitempty"`
		Expiration            *time.Time  `json:"expiration,omitempty"`
	}

	DelegatesWrapper struct {
		DaoID          string             `json:"dao_id"`
		DelegationType string             `json:"delegation_type"`
		ChainId        *string            `json:"chain_id,omitempty"`
		TotalCnt       int32              `json:"total_cnt"`
		Delegates      []*DelegateEntryV2 `json:"list"`
	}

	GetDelegatesV2Response struct {
		List     []*DelegatesWrapper `json:"list"`
		TotalCnt int32               `json:"total_cnt"`
	}

	GetDelegatorsV2Response struct {
		List     []*DelegatesWrapper `json:"list"`
		TotalCnt int32               `json:"total_cnt"`
	}

	GetUserDelegatesV2Response struct {
		List     []*DelegatesWrapper `json:"list"`
		TotalCnt int32               `json:"total_cnt"`
	}

	GetUserDelegatesTopV2Response struct {
		List     []*DelegatesWrapper `json:"list"`
		TotalCnt int32               `json:"total_cnt"`
	}

	GetUserDelegatorsV2Response struct {
		List     []*DelegatesWrapper `json:"list"`
		TotalCnt int32               `json:"total_cnt"`
	}

	GetUserDelegatorsTopV2Response struct {
		List     []*DelegatesWrapper `json:"list"`
		TotalCnt int32               `json:"total_cnt"`
	}
)
