package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go.openly.dev/pointy"

	"github.com/gorilla/mux"
	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
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

func (h *Delegate) EnrichRoutes(v1, v2 *mux.Router) {
	v1.HandleFunc("/user/{address}/delegates/top", h.getDelegatesByAddress).Methods(http.MethodGet).Name("get_delegates_by_address")
	v1.HandleFunc("/user/{address}/delegators/top", h.getDelegatorsByAddress).Methods(http.MethodGet).Name("get_delegators_by_address")
	v1.HandleFunc("/user/{address}/delegations/total", h.getTotalDelegations).Methods(http.MethodGet).Name("get_delegates_summary_by_address")
	v1.HandleFunc("/user/{address}/delegates/{dao_id}/list", h.getDelegatesList).Methods(http.MethodGet).Name("get_delegates_list")
	v1.HandleFunc("/user/{address}/delegators/{dao_id}/list", h.getDelegatorsList).Methods(http.MethodGet).Name("get_delegators_list")

	v2.HandleFunc("/user/{address}/delegates/top", h.getUserDelegatesTopV2).Methods(http.MethodGet).Name("get_user_delegates_top_v2")
	v2.HandleFunc("/user/{address}/delegators/top", h.getUserDelegatorsTopV2).Methods(http.MethodGet).Name("get_user_delegates_top_v2")
	v2.HandleFunc("/user/{address}/delegates/{dao_id}/list", h.getUserDelegatesListV2).Methods(http.MethodGet).Name("get_delegates_list")
	v2.HandleFunc("/user/{address}/delegators/{dao_id}/list", h.getUserDelegatorsListV2).Methods(http.MethodGet).Name("get_delegators_list")
}

func (h *Delegate) getDelegatesByAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	resp, err := h.dc.GetTopDelegates(r.Context(), &storagepb.GetTopDelegatesRequest{Address: address})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get delegates by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convertToTopDelegatesFromProto(resp))
}

func (h *Delegate) getDelegatorsByAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	resp, err := h.dc.GetTopDelegators(r.Context(), &storagepb.GetTopDelegatorsRequest{Address: address})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get delegators by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convertToTopDelegatorsFromProto(resp))
}

func (h *Delegate) getTotalDelegations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	resp, err := h.dc.GetDelegationSummary(r.Context(), &storagepb.GetDelegationSummaryRequest{Address: address})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get GetDelegatesSummary by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(delegate.TotalDelegations{
		TotalDelegatorsCount: int(resp.GetTotalDelegatorsCount()),
		TotalDelegatesCount:  int(resp.GetTotalDelegatesCount()),
	})
}

func convertToTopDelegatesFromProto(info *storagepb.GetTopDelegatesResponse) delegate.TopDelegates {
	all := delegate.TopDelegates{
		TotalCount: int(info.GetTotalDelegatesCount()),
		List:       make([]delegate.DelegationSummary, 0, len(info.GetList())),
	}

	for _, delegation := range info.GetList() {
		daoDelegations := make([]delegate.DelegationDetails, 0, len(delegation.GetList()))
		for _, dd := range delegation.GetList() {
			daoDelegations = append(daoDelegations, convertDelegationToModel(dd))
		}

		all.List = append(all.List, delegate.DelegationSummary{
			Dao:        convertToDaoFromProto(delegation.Dao),
			List:       daoDelegations,
			TotalCount: int(delegation.TotalCount),
		})
	}

	return all
}

func convertToTopDelegatorsFromProto(info *storagepb.GetTopDelegatorsResponse) delegate.TopDelegators {
	all := delegate.TopDelegators{
		TotalCount: int(info.GetTotalDelegatorsCount()),
		List:       make([]delegate.DelegationSummary, 0, len(info.GetList())),
	}

	for _, di := range info.GetList() {
		daoDelegations := make([]delegate.DelegationDetails, 0, len(di.GetList()))
		for _, dd := range di.GetList() {
			daoDelegations = append(daoDelegations, convertDelegationToModel(dd))
		}

		all.List = append(all.List, delegate.DelegationSummary{
			Dao:        convertToDaoFromProto(di.Dao),
			List:       daoDelegations,
			TotalCount: int(di.TotalCount),
		})
	}

	return all
}

func (h *Delegate) getDelegatesList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	daoID := vars["dao_id"]

	form, verr := common.NewPagination().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)
		return
	}

	params := form.(*common.Pagination)
	resp, err := h.dc.GetDelegatesByDao(r.Context(), &storagepb.GetDelegatesByDaoRequest{
		DaoId:   daoID,
		Address: address,
		Limit:   uint32(params.Limit),
		Offset:  pointy.Uint32(uint32(params.Offset)),
	})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get delegates by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convertToDelegatesListFromProto(resp))
}

func convertToDelegatesListFromProto(info *storagepb.GetDelegatesByDaoResponse) delegate.DelegatesList {
	all := delegate.DelegatesList{
		List:       make([]delegate.DelegationDetails, 0, len(info.GetList())),
		TotalCount: int(info.GetTotalCount()),
	}

	for _, di := range info.GetList() {
		all.List = append(all.List, convertDelegationToModel(di))
	}

	return all
}

func convertToDelegatorsListFromProto(info *storagepb.GetDelegatorsByDaoResponse) delegate.DelegatorsList {
	all := delegate.DelegatorsList{
		List:       make([]delegate.DelegationDetails, 0, len(info.GetList())),
		TotalCount: int(info.GetTotalCount()),
	}

	for _, di := range info.GetList() {
		all.List = append(all.List, convertDelegationToModel(di))
	}

	return all
}

func convertDelegationToModel(info *storagepb.DelegationDetails) delegate.DelegationDetails {
	var exp *time.Time
	if info.GetExpiration() != nil {
		expTime := info.GetExpiration().AsTime()

		// small fix for getting correct expirations
		if expTime.Before(time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local)) {
			exp = &expTime
		}
	}

	return delegate.DelegationDetails{
		Address:             info.GetAddress(),
		EnsName:             info.GetEnsName(),
		PercentOfDelegators: int(info.GetPercentOfDelegators()),
		Expiration:          exp,
	}
}

func (h *Delegate) getDelegatorsList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	daoID := vars["dao_id"]

	form, verr := common.NewPagination().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)
		return
	}

	params := form.(*common.Pagination)
	resp, err := h.dc.GetDelegatorsByDao(r.Context(), &storagepb.GetDelegatorsByDaoRequest{
		DaoId:   daoID,
		Address: address,
		Limit:   uint32(params.Limit),
		Offset:  pointy.Uint32(uint32(params.Offset)),
	})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get delegates by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convertToDelegatorsListFromProto(resp))
}
