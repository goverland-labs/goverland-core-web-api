package dao

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
}
