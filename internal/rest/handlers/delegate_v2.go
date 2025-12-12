package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	proto "github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/rs/zerolog/log"
	"go.openly.dev/pointy"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	forms "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/dao"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/delegate"
)

func (h *DAO) getDelegatesV2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	daoID := vars["id"]

	form, verr := forms.NewGetDelegatesV2Form().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetDelegatesV2)

	// TODO: for now we only support one address in query
	var qAccounts []string
	if params.Query != nil {
		qAccounts = append(qAccounts, *params.Query)
	}

	resp, err := h.delegateClient.GetDelegatesV2(r.Context(), &proto.GetDelegatesV2Request{
		DaoId:          daoID,
		QueryAccounts:  qAccounts,
		Sort:           params.By,
		Limit:          int32(params.Limit),
		Offset:         int32(params.Offset),
		DelegationType: convertDelegationTypeToProto(pointy.StringValue(params.DelegationType, "")),
		ChainId:        params.ChainID,
	})
	if err != nil {
		log.Error().Err(err).Msg("get dao delegates v2")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	list := make([]*delegate.DelegatesWrapper, 0, len(resp.List))
	for _, v := range resp.List {
		list = append(list, convertDelegateWrapperToModel(v))
	}

	result := delegate.GetDelegatesV2Response{
		List:     list,
		TotalCnt: resp.TotalCnt,
	}

	_ = json.NewEncoder(w).Encode(result)
}

func convertDelegateWrapperToModel(info *proto.DelegatesWrapper) *delegate.DelegatesWrapper {
	if info == nil {
		return nil
	}

	delegates := make([]*delegate.DelegateEntryV2, 0, len(info.Delegates))
	for _, di := range info.Delegates {
		delegates = append(delegates, convertDelegateEntryV2ToModel(di))
	}

	return &delegate.DelegatesWrapper{
		DaoID:          info.GetDaoId(),
		DelegationType: convertDelegationTypeToInternal(info.GetDelegationType()),
		ChainId:        info.ChainId,
		TotalCnt:       info.GetTotalCnt(),
		Delegates:      delegates,
	}
}

func convertDelegateEntryV2ToModel(entry *proto.DelegateEntryV2) *delegate.DelegateEntryV2 {
	if entry == nil {
		return nil
	}

	var expiration *time.Time
	if entry.Expiration != nil && !entry.Expiration.AsTime().IsZero() {
		exp := entry.GetExpiration().AsTime()
		expiration = &exp
	}

	var tv *delegate.TokenValue
	if entry.TokenValue != nil {
		tv = &delegate.TokenValue{
			Value:    entry.GetTokenValue().GetValue(),
			Symbol:   entry.GetTokenValue().GetSymbol(),
			Decimals: entry.GetTokenValue().GetDecimals(),
		}
	}

	return &delegate.DelegateEntryV2{
		Address:               entry.GetAddress(),
		EnsName:               entry.GetEnsName(),
		DelegatorCount:        entry.DelegatorCount,
		PercentOfDelegators:   entry.PercentOfDelegators,
		PercentOfVotingPower:  entry.PercentOfVotingPower,
		About:                 entry.About,
		Statement:             entry.Statement,
		VotesCount:            entry.VotesCount,
		CreatedProposalsCount: entry.CreatedProposalsCount,
		VotingPower:           entry.VotingPower,
		TokenValue:            tv,
		Expiration:            expiration,
	}
}

func (h *DAO) getUserDelegatorsV2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	daoID := vars["id"]
	address := vars["address"]

	form, verr := forms.NewGetDelegatorsV2Form().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetDelegatorsV2)
	resp, err := h.delegateClient.GetDelegatorsV2(r.Context(), &proto.GetDelegatorsV2Request{
		DaoId:          daoID,
		QueryAccounts:  []string{address},
		Limit:          int32(params.Limit),
		Offset:         int32(params.Offset),
		DelegationType: convertDelegationTypeToProto(pointy.StringValue(params.DelegationType, "")),
		ChainId:        params.ChainID,
	})
	if err != nil {
		log.Error().Err(err).Msg("get dao user delegators v2")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	list := make([]*delegate.DelegatesWrapper, 0, len(resp.List))
	for _, v := range resp.List {
		list = append(list, convertDelegateWrapperToModel(v))
	}

	result := delegate.GetDelegatorsV2Response{
		List:     list,
		TotalCnt: resp.TotalCnt,
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *DAO) getUserDelegatorsTopV2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	daoID := vars["id"]
	address := vars["address"]

	resp, err := h.delegateClient.GetTopDelegatorsV2(r.Context(), &proto.GetTopDelegatorsV2Request{
		DaoId:   &daoID,
		Address: address,
	})
	if err != nil {
		log.Error().Err(err).Msg("get dao user top delegators v2")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	list := make([]*delegate.DelegatesWrapper, 0, len(resp.List))
	for _, v := range resp.List {
		list = append(list, convertDelegateWrapperToModel(v))
	}

	result := delegate.GetDelegatorsV2Response{
		List:     list,
		TotalCnt: resp.TotalCnt,
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *Delegate) getUserDelegatesTopV2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	resp, err := h.dc.GetTopDelegatesV2(r.Context(), &proto.GetTopDelegatesV2Request{Address: address})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get top delegates v2 by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	list := make([]*delegate.DelegatesWrapper, 0, len(resp.List))
	for _, v := range resp.List {
		list = append(list, convertDelegateWrapperToModel(v))
	}

	result := delegate.GetUserDelegatesTopV2Response{
		List:     list,
		TotalCnt: resp.TotalCnt,
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *Delegate) getUserDelegatesListV2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	daoID := vars["dao_id"]
	address := vars["address"]

	form, verr := forms.NewGetUserDelegatesV2Form().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetUserDelegatesV2)

	resp, err := h.dc.GetDelegatesV2(r.Context(), &proto.GetDelegatesV2Request{
		DaoId:          daoID,
		QueryAccounts:  []string{address},
		Limit:          int32(params.Limit),
		Offset:         int32(params.Offset),
		DelegationType: convertDelegationTypeToProto(pointy.StringValue(params.DelegationType, "")),
		ChainId:        params.ChainID,
	})
	if err != nil {
		log.Error().Err(err).Msg("get dao delegates v2")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	list := make([]*delegate.DelegatesWrapper, 0, len(resp.List))
	for _, v := range resp.List {
		list = append(list, convertDelegateWrapperToModel(v))
	}

	result := delegate.GetUserDelegatesV2Response{
		List:     list,
		TotalCnt: resp.TotalCnt,
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *Delegate) getUserDelegatorsTopV2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	resp, err := h.dc.GetTopDelegatorsV2(r.Context(), &proto.GetTopDelegatorsV2Request{Address: address})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"address": address,
		}).Msg("get top delegators v2 by address")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	list := make([]*delegate.DelegatesWrapper, 0, len(resp.List))
	for _, v := range resp.List {
		list = append(list, convertDelegateWrapperToModel(v))
	}

	result := delegate.GetUserDelegatorsTopV2Response{
		List:     list,
		TotalCnt: resp.TotalCnt,
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *Delegate) getUserDelegatorsListV2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	daoID := vars["dao_id"]
	address := vars["address"]

	form, verr := forms.NewGetUserDelegatorsV2Form().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetUserDelegatorsV2)

	resp, err := h.dc.GetDelegatorsV2(r.Context(), &proto.GetDelegatorsV2Request{
		DaoId:          daoID,
		QueryAccounts:  []string{address},
		Limit:          int32(params.Limit),
		Offset:         int32(params.Offset),
		DelegationType: convertDelegationTypeToProto(pointy.StringValue(params.DelegationType, "")),
		ChainId:        params.ChainID,
	})
	if err != nil {
		log.Error().Err(err).Msg("get dao delegators v2")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	list := make([]*delegate.DelegatesWrapper, 0, len(resp.List))
	for _, v := range resp.List {
		list = append(list, convertDelegateWrapperToModel(v))
	}

	result := delegate.GetUserDelegatorsV2Response{
		List:     list,
		TotalCnt: resp.TotalCnt,
	}

	_ = json.NewEncoder(w).Encode(result)
}
