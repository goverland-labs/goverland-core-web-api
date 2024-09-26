package delegate

import (
	"time"

	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/dao"
)

type DelegationDetails struct {
	// The delegation address
	Address string `json:"address"`
	// Resolved ens name
	EnsName string `json:"ens_name,omitempty"`
	// Percentage of delegation
	PercentOfDelegators int `json:"percent_of_delegators,omitempty"`
	// Expires at date. If 0 the expiration is not set
	Expiration *time.Time `json:"expiration,omitempty"`
}

type DelegationSummary struct {
	// Dao details
	Dao dao.Dao `json:"dao"`
	// List of delegations
	Delegations []DelegationDetails `json:"delegations,omitempty"`
}

type AllDelegations struct {
	// The number of total delegations in out DB
	TotalDelegationsCount int `json:"total_delegations_count"`
	// List of delegations grouped by dao and sorted by popularity index
	Delegations []DelegationSummary `json:"delegations,omitempty"`
}
