package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/goverland-labs/core-api/protobuf/internalapi"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/core-web-api/internal/response"
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
	baseRouter.HandleFunc("/dao/{id}", h.getByIDAction).Methods(http.MethodGet).Name("get_dao_by_id")
}

func (h *DAO) getByIDAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dao, err := h.dc.GetByID(r.Context(), &internalapi.DaoByIDRequest{DaoId: id})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"id": id,
		}).Msg("get dao by id")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convertToDaoFromProto(dao.Dao))
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
