package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/goverland-labs/core-api/protobuf/internalapi"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/core-web-api/internal/response"
	forms "github.com/goverland-labs/core-web-api/internal/rest/form/dao"
	"github.com/goverland-labs/core-web-api/internal/rest/models/dao"
)

type DAO struct {
	dc internalapi.DaoClient
}

func NewDaoHandler(dc internalapi.DaoClient) APIHandler {
	return &DAO{
		dc: dc,
	}
}

func (h *DAO) EnrichRoutes(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/dao/top", h.getTopAction).Methods(http.MethodGet).Name("get_dao_top")
	baseRouter.HandleFunc("/dao/{id}", h.getByIDAction).Methods(http.MethodGet).Name("get_dao_by_id")
	baseRouter.HandleFunc("/daos", h.getListAction).Methods(http.MethodGet).Name("get_dao_list")
}

func (h *DAO) getByIDAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.dc.GetByID(r.Context(), &internalapi.DaoByIDRequest{DaoId: id})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"id": id,
		}).Msg("get dao by id")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convertToDaoFromProto(resp.Dao))
}

func (h *DAO) getListAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewGetListForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetList)
	list, err := h.dc.GetByFilter(r.Context(), &internalapi.DaoByFilterRequest{
		Query:    &params.Query,
		Category: &params.Category,
		Limit:    &params.Limit,
		Offset:   &params.Offset,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("get dao list by filter")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	resp := make([]dao.Dao, len(list.Daos))
	for i, info := range list.Daos {
		resp[i] = convertToDaoFromProto(info)
	}

	response.AddPaginationHeaders(w, params.Offset, params.Limit, list.TotalCount)

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *DAO) getTopAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewGetTopForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetTop)
	list, err := h.dc.GetTopByCategories(r.Context(), &internalapi.TopByCategoriesRequest{
		Limit: params.Limit,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("get top dao")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	resp := make(dao.TopCategories)
	for _, info := range list.GetCategories() {
		daos := make([]dao.Dao, len(info.GetDaos()))
		for i, details := range info.GetDaos() {
			daos[i] = convertToDaoFromProto(details)
		}

		resp[info.GetCategory()] = dao.TopCategory{
			TotalCount: info.GetTotalCount(),
			List:       daos,
		}
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func convertToDaoFromProto(info *internalapi.DaoInfo) dao.Dao {
	return dao.Dao{
		ID:             info.GetId(),
		CreatedAt:      info.GetCreatedAt().AsTime(),
		UpdatedAt:      info.GetUpdatedAt().AsTime(),
		Name:           info.GetName(),
		Private:        info.GetPrivate(),
		About:          info.GetAbout(),
		Avatar:         info.GetAvatar(),
		Terms:          info.GetTerms(),
		Location:       info.GetLocation(),
		Website:        info.GetWebsite(),
		Twitter:        info.GetTwitter(),
		Github:         info.GetGithub(),
		Coingecko:      info.GetCoingeko(),
		Email:          info.GetEmail(),
		Network:        info.GetNetwork(),
		Symbol:         info.GetSymbol(),
		Skin:           info.GetSkin(),
		Domain:         info.GetDomain(),
		Strategies:     convertToStrategiesFromProto(info.GetStrategies()),
		Voting:         convertToVotingFromProto(info.GetVoting()),
		Categories:     info.GetCategories(),
		Treasures:      convertToTreasuresFromProto(info.GetTreasuries()),
		FollowersCount: info.GetFollowersCount(),
		ProposalsCount: info.GetProposalsCount(),
		Guidelines:     info.GetGuidelines(),
		Template:       info.GetTemplate(),
		ParentID:       info.GetParentId(),
	}
}

func convertToStrategiesFromProto(info []*internalapi.Strategy) dao.Strategies {
	res := make(dao.Strategies, len(info))

	for i, details := range info {
		res[i] = dao.Strategy{
			Name:    details.GetName(),
			Network: details.GetNetwork(),
		}
	}

	return res
}

func convertToTreasuresFromProto(info []*internalapi.Treasury) dao.Treasuries {
	res := make(dao.Treasuries, len(info))

	for i, details := range info {
		res[i] = dao.Treasury{
			Name:    details.GetName(),
			Address: details.GetAddress(),
			Network: details.GetNetwork(),
		}
	}

	return res
}

func convertToVotingFromProto(info *internalapi.Voting) dao.Voting {
	return dao.Voting{
		Delay:       info.GetDelay(),
		Period:      info.GetPeriod(),
		Type:        info.GetType(),
		Quorum:      info.GetQuorum(),
		Blind:       info.GetBlind(),
		HideAbstain: info.GetHideAbstain(),
		Privacy:     info.GetPrivacy(),
		Aliased:     info.GetAliased(),
	}
}
