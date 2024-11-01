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
	List []DelegationDetails `json:"list,omitempty"`
	// The number of delegations for DAO
	TotalCount int `json:"total_count"`
}

type TopDelegates struct {
	// The number of total delegations in out DB
	TotalCount int `json:"total_count"`
	// List of delegations grouped by dao and sorted by popularity index
	List []DelegationSummary `json:"list,omitempty"`
}

type TopDelegators struct {
	// The number of total delegators in out DB
	TotalCount int `json:"total_count"`
	// List of delegators grouped by dao and sorted by popularity index
	List []DelegationSummary `json:"list,omitempty"`
}

type TotalDelegations struct {
	// The number of total delegators in out DB
	TotalDelegatorsCount int `json:"total_delegators_count"`
	// The number of total delegates in out DB
	TotalDelegatesCount int `json:"total_delegates_count"`
}

type DelegatesList struct {
	// List of delegations
	List []DelegationDetails `json:"list,omitempty"`
	// The number of delegations for DAO
	TotalCount int `json:"total_count"`
}

type DelegatorsList struct {
	// List of delegations
	List []DelegationDetails `json:"list,omitempty"`
	// The number of delegations for DAO
	TotalCount int `json:"total_count"`
}
