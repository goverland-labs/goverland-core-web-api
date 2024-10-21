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
	baseRouter.HandleFunc("/delegates", h.getDelegatesByAddress).Methods(http.MethodGet).Name("get_delegates_by_address")
	baseRouter.HandleFunc("/delegators", h.getDelegatorsByAddress).Methods(http.MethodGet).Name("get_delegators_by_address")
	baseRouter.HandleFunc("/delegations/total", h.getTotalDelegations).Methods(http.MethodGet).Name("get_delegates_summary_by_address")
}

func (h *Delegate) getDelegatesByAddress(w http.ResponseWriter, r *http.Request) {
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

func (h *Delegate) getDelegatorsByAddress(w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")

	resp, err := h.dc.GetAllDelegators(r.Context(), &storagepb.GetAllDelegatorsRequest{Address: address})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get delegators by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convertToDelegatorsFromProto(resp))
}

func (h *Delegate) getTotalDelegations(w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")

	resp, err := h.dc.GetDelegatesSummary(r.Context(), &storagepb.GetDelegatesSummaryRequest{Address: address})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get delegators by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(delegate.TotalDelegations{
		TotalDelegatorsCount:  int(resp.GetTotalDelegatorsCount()),
		TotalDelegationsCount: int(resp.GetTotalDelegationsCount()),
	})
}

func convertToDelegatesFromProto(info *storagepb.GetAllDelegationsResponse) delegate.AllDelegations {
	all := delegate.AllDelegations{
		TotalDelegationsCount: int(info.GetTotalDelegationsCount()),
		Delegations:           make([]delegate.DelegationSummary, 0, len(info.GetDelegations())),
	}

	for _, delegation := range info.GetDelegations() {
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

func convertToDelegatorsFromProto(info *storagepb.GetAllDelegatorsResponse) delegate.AllDelegators {
	all := delegate.AllDelegators{
		TotalDelegatorsCount: int(info.GetTotalDelegatorsCount()),
		Delegations:          make([]delegate.DelegationSummary, 0, len(info.GetDelegators())),
	}

	for _, di := range info.GetDelegators() {
		daoDelegations := make([]delegate.DelegationDetails, 0, len(di.GetDelegators()))
		for _, dd := range di.GetDelegators() {
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
			Dao:         convertToDaoFromProto(di.Dao),
			Delegations: daoDelegations,
		})
	}

	return all
}
