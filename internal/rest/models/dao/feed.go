package dao

import (
	"encoding/json"
	"time"
)

type FeedItem struct {
	ID           uint64          `json:"id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DaoID        string          `json:"dao_id"`
	ProposalID   string          `json:"proposal_id"`
	DiscussionID string          `json:"discussion_id"`
	Type         string          `json:"type"`
	Action       string          `json:"action"`
	Snapshot     json.RawMessage `json:"snapshot"`
}
