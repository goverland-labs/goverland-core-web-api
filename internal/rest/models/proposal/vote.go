package proposal

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Vote struct {
	ID           string          `json:"id"`
	Ipfs         string          `json:"ipfs"`
	DaoID        uuid.UUID       `json:"dao_id"`
	ProposalID   string          `json:"proposal_id"`
	Voter        string          `json:"voter"`
	EnsName      string          `json:"ens_name"`
	Created      uint64          `json:"created"`
	Reason       string          `json:"reason"`
	Choice       json.RawMessage `json:"choice"`
	App          string          `json:"app"`
	Vp           float32         `json:"vp"`
	VpByStrategy []float32       `json:"vp_by_strategy"`
	VpState      string          `json:"vp_state"`
}

type VoteValidation struct {
	OK                  bool                 `json:"ok"`
	VotingPower         float64              `json:"voting_power"`
	VoteValidationError *VoteValidationError `json:"error,omitempty"`
	VoteStatus          VoteStatus           `json:"status"`
}

type VoteStatus struct {
	Voted  bool            `json:"voted"`
	Choice json.RawMessage `json:"choice,omitempty"`
}

type VoteValidationError struct {
	Message string `json:"message"`
	Code    uint32 `json:"code"`
}

type VotePreparation struct {
	ID        string `json:"id"`
	TypedData string `json:"typed_data"`
}

type SuccessfulVote struct {
	ID         string  `json:"id"`
	IPFS       string  `json:"ipfs"`
	Relayer    Relayer `json:"relayer"`
	ProposalID string  `json:"proposal_id"`
}

type Relayer struct {
	Address string `json:"address"`
	Receipt string `json:"receipt"`
}
