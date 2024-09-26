package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/delegate"
)

type Delegate struct {
	dc storagepb.DelegateClient
}

func NewDelegateHandler(dc storagepb.DelegateClient) APIHandler {
	return &Delegate{
		dc: dc,
	}
}

func (h *Delegate) EnrichRoutes(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/delegates", h.getByAddress).Methods(http.MethodGet).Name("get_delegates_by_address")
}

func (h *Delegate) getByAddress(w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")

	resp, err := h.dc.GetAllDelegations(r.Context(), &storagepb.GetAllDelegationsRequest{Address: address})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get delegates by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convertToDelegatesFromProto(resp))
}

func convertToDelegatesFromProto(info *storagepb.GetAllDelegationsResponse) delegate.AllDelegations {
	all := delegate.AllDelegations{
		TotalDelegationsCount: int(info.GetTotalDelegationsCount()),
		Delegations:           make([]delegate.DelegationSummary, 0, len(info.Delegations)),
	}

	for _, delegation := range info.Delegations {
		daoDelegations := make([]delegate.DelegationDetails, 0, len(delegation.Delegations))
		for _, dd := range delegation.Delegations {
			var exp *time.Time
			if dd.GetExpiration() != nil {
				expTime := dd.GetExpiration().AsTime()
				exp = &expTime
			}

			daoDelegations = append(daoDelegations, delegate.DelegationDetails{
				Address:             dd.GetAddress(),
				EnsName:             dd.GetEnsName(),
				PercentOfDelegators: int(dd.GetPercentOfDelegators()),
				Expiration:          exp,
			})
		}

		all.Delegations = append(all.Delegations, delegate.DelegationSummary{
			Dao:         convertToDaoFromProto(delegation.Dao),
			Delegations: daoDelegations,
		})
	}

	return all
}
