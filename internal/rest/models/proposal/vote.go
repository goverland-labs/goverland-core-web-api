package proposal

import (
	"github.com/google/uuid"
)

type Vote struct {
	ID           string    `json:"id"`
	Ipfs         string    `json:"ipfs"`
	DaoID        uuid.UUID `json:"dao_id"`
	ProposalID   string    `json:"proposal_id"`
	Voter        string    `json:"voter"`
	Created      uint64    `json:"created"`
	Reason       string    `json:"reason"`
	Choice       int       `json:"choice"`
	App          string    `json:"app"`
	Vp           float32   `json:"vp"`
	VpByStrategy []float32 `json:"vp_by_strategy"`
	VpState      string    `json:"vp_state"`
}
