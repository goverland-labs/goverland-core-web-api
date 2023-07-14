package dao

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	None                        TimelineAction = ""
	DaoCreated                  TimelineAction = "dao.created"
	DaoUpdated                  TimelineAction = "dao.updated"
	ProposalCreated             TimelineAction = "proposal.created"
	ProposalUpdated             TimelineAction = "proposal.updated"
	ProposalVotingStartsSoon    TimelineAction = "proposal.voting.starts_soon"
	ProposalVotingStarted       TimelineAction = "proposal.voting.started"
	ProposalVotingQuorumReached TimelineAction = "proposal.voting.quorum_reached"
	ProposalVotingEnded         TimelineAction = "proposal.voting.ended"
)

type FeedItem struct {
	ID           uuid.UUID       `json:"id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DaoID        uuid.UUID       `json:"dao_id"`
	ProposalID   string          `json:"proposal_id"`
	DiscussionID string          `json:"discussion_id"`
	Type         string          `json:"type"`
	Action       string          `json:"action"`
	Snapshot     json.RawMessage `json:"snapshot"`
	Timeline     []TimelineItem  `json:"timeline"`
}

type TimelineItem struct {
	CreatedAt time.Time `json:"created_at"`
	Action    TimelineAction
}

type TimelineAction string
