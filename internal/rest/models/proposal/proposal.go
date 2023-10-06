package proposal

import (
	"time"

	"github.com/google/uuid"
)

type Choices []string

type Scores []float32

type Strategy struct {
	Name    string                 `json:"name"`
	Network string                 `json:"network"`
	Params  map[string]interface{} `json:"params"`
}

type Strategies []Strategy

const (
	None                TimelineAction = ""
	Created             TimelineAction = "proposal.created"
	Updated             TimelineAction = "proposal.updated"
	VotingStartsSoon    TimelineAction = "proposal.voting.starts_soon"
	VotingEndsSoon      TimelineAction = "proposal.voting.ends_soon"
	VotingStarted       TimelineAction = "proposal.voting.started"
	VotingQuorumReached TimelineAction = "proposal.voting.quorum_reached"
	VotingEnded         TimelineAction = "proposal.voting.ended"
)

type TimelineAction string

type TimelineItem struct {
	CreatedAt time.Time      `json:"created_at"`
	Action    TimelineAction `json:"action"`
}

type Proposal struct {
	ID            string         `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Ipfs          string         `json:"ipfs"`
	Author        string         `json:"author"`
	Created       uint64         `json:"created"`
	DaoID         uuid.UUID      `json:"dao_id"`
	Network       string         `json:"network"`
	Symbol        string         `json:"symbol"`
	Type          string         `json:"type"`
	Strategies    Strategies     `json:"strategies"`
	Title         string         `json:"title"`
	Body          string         `json:"body"`
	Discussion    string         `json:"discussion"`
	Choices       Choices        `json:"choices"`
	Start         uint64         `json:"start"`
	End           uint64         `json:"end"`
	Quorum        float32        `json:"quorum"`
	Privacy       string         `json:"privacy"`
	Snapshot      string         `json:"snapshot"`
	State         string         `json:"state"`
	Link          string         `json:"link"`
	App           string         `json:"app"`
	Scores        Scores         `json:"scores"`
	ScoresState   string         `json:"scores_state"`
	ScoresTotal   float32        `json:"scores_total"`
	ScoresUpdated uint64         `json:"scores_updated"`
	Votes         uint64         `json:"votes"`
	Timeline      []TimelineItem `json:"timeline,omitempty"`
}
