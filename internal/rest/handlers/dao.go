package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/goverland-labs/goverland-core-feed/protocol/feedpb"
	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/rs/zerolog/log"
	"go.openly.dev/pointy"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	forms "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/dao"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/dao"
	"github.com/goverland-labs/goverland-core-web-api/pkg/helpers"
)

const (
	delegationTypeSplitDelegation = "split-delegation"
	delegationTypeDelegation      = "delegation"
	delegationTypeErc20Votes      = "erc20-votes"
)

type DAO struct {
	dc             storagepb.DaoClient
	fc             feedpb.FeedClient
	delegateClient storagepb.DelegateClient
}

func NewDaoHandler(dc storagepb.DaoClient, fc feedpb.FeedClient, delegateClient storagepb.DelegateClient) APIHandler {
	return &DAO{
		dc:             dc,
		fc:             fc,
		delegateClient: delegateClient,
	}
}

func (h *DAO) EnrichRoutes(v1, v2 *mux.Router) {
	v1.HandleFunc("/daos/top", h.getTopAction).Methods(http.MethodGet).Name("get_dao_top")
	v1.HandleFunc("/daos/recommendations", h.getRecommendations).Methods(http.MethodGet).Name("get_dao_recommendations")
	v1.HandleFunc("/daos/{id}/feed", h.getFeedByIDAction).Methods(http.MethodGet).Name("get_dao_feed_by_id")
	v1.HandleFunc("/daos/{id}", h.getByIDAction).Methods(http.MethodGet).Name("get_dao_by_id")
	v1.HandleFunc("/daos", h.getListAction).Methods(http.MethodGet).Name("get_dao_list")
	v1.HandleFunc("/daos/{id}/delegates", h.getDelegates).Methods(http.MethodGet).Name("get_delegates_list")
	v1.HandleFunc("/daos/{id}/delegate-profile", h.getDelegateProfile).Methods(http.MethodGet).Name("get_delegate_profile")
	v1.HandleFunc("/daos/{id}/delegates/{address}/delegators", h.getDelegators).Methods(http.MethodGet).Name("get_delegators")
	v1.HandleFunc("/daos/{id}/token-info", h.getTokenInfo).Methods(http.MethodGet).Name("get_dao_token_info")
	v1.HandleFunc("/daos/{id}/token-chart", h.getTokenChart).Methods(http.MethodGet).Name("get_dao_token_chart")
	v1.HandleFunc("/daos/{id}/populate-token-price", h.populateTokenPrice).Methods(http.MethodPost).Name("populate_dao_token_price")
	v1.HandleFunc("/daos/update-fungible-ids", h.updateFungibleIds).Methods(http.MethodPost).Name("update_fungible_ids")

	v2.HandleFunc("/daos/{id}/delegates", h.getDelegatesV2).Methods(http.MethodGet).Name("get_delegates_v2_list")
	v2.HandleFunc("/daos/{id}/delegates/{address}/delegators", h.getUserDelegatorsV2).Methods(http.MethodGet).Name("get_delegators_v2_list")
}

func (h *DAO) getByIDAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.dc.GetByID(r.Context(), &storagepb.DaoByIDRequest{DaoId: id})
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

func (h *DAO) getFeedByIDAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	form, verr := forms.NewGetFeedForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetFeed)
	resp, err := h.fc.GetByFilter(r.Context(), &feedpb.FeedByFilterRequest{
		DaoId:  &id,
		Types:  []string{"proposal"}, // todo: move to api
		Limit:  &params.Limit,
		Offset: &params.Offset,
	})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"id": id,
		}).Msg("get feed by dao")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	list := make([]dao.FeedItem, len(resp.Items))
	for i, fi := range resp.Items {
		list[i] = convertToFeedItemFromProto(fi)
	}

	response.AddPaginationHeaders(w, params.Offset, params.Limit, resp.TotalCount)
	_ = json.NewEncoder(w).Encode(list)
}

func convertToFeedItemFromProto(fi *feedpb.FeedInfo) dao.FeedItem {
	itemID, _ := uuid.Parse(fi.GetId())
	daoID, _ := uuid.Parse(fi.GetDaoId())

	return dao.FeedItem{
		ID:           itemID,
		CreatedAt:    fi.GetCreatedAt().AsTime(),
		UpdatedAt:    fi.GetUpdatedAt().AsTime(),
		DaoID:        daoID,
		ProposalID:   fi.GetProposalId(),
		DiscussionID: fi.GetDiscussionId(),
		Type:         convertProtoType(fi.GetType()),
		Action:       fi.GetAction(),
		Snapshot:     fi.GetSnapshot().Value,
		Timeline:     convertTimelineFromProto(fi.GetTimeline()),
	}
}

// todo: move to constant
func convertProtoType(ft feedpb.FeedInfo_Type) string {
	switch ft {
	case feedpb.FeedInfo_Proposal:
		return "proposal"
	case feedpb.FeedInfo_DAO:
		return "dao"
	default:
		return "unspecified"
	}
}

func convertTimelineFromProto(timeline []*feedpb.FeedTimelineItem) []dao.TimelineItem {
	converted := make([]dao.TimelineItem, 0, len(timeline))

	for _, t := range timeline {
		converted = append(converted, dao.TimelineItem{
			CreatedAt: t.GetCreatedAt().AsTime(),
			Action:    convertTimelineActionProto(t.GetAction()),
		})
	}

	return converted
}

var timelineActionMap = map[feedpb.FeedTimelineItem_TimelineAction]dao.TimelineAction{
	feedpb.FeedTimelineItem_DaoCreated:                  dao.DaoCreated,
	feedpb.FeedTimelineItem_DaoUpdated:                  dao.DaoUpdated,
	feedpb.FeedTimelineItem_ProposalCreated:             dao.ProposalCreated,
	feedpb.FeedTimelineItem_ProposalUpdated:             dao.ProposalUpdated,
	feedpb.FeedTimelineItem_ProposalVotingStartsSoon:    dao.ProposalVotingStartsSoon,
	feedpb.FeedTimelineItem_ProposalVotingStarted:       dao.ProposalVotingStarted,
	feedpb.FeedTimelineItem_ProposalVotingQuorumReached: dao.ProposalVotingQuorumReached,
	feedpb.FeedTimelineItem_ProposalVotingEnded:         dao.ProposalVotingEnded,
}

func convertTimelineActionProto(action feedpb.FeedTimelineItem_TimelineAction) dao.TimelineAction {
	converted, exists := timelineActionMap[action]
	if !exists {
		log.Warn().Str("action", action.String()).Msg("unknown timeline action")
		return dao.None
	}

	return converted
}

func (h *DAO) getListAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewGetListForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetList)
	list, err := h.dc.GetByFilter(r.Context(), &storagepb.DaoByFilterRequest{
		Query:       params.Query,
		Category:    params.Category,
		Limit:       &params.Limit,
		Offset:      &params.Offset,
		DaoIds:      params.DAOs,
		FungibleIds: params.FungibleIDs,
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
	list, err := h.dc.GetTopByCategories(r.Context(), &storagepb.TopByCategoriesRequest{
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

func (h *DAO) getRecommendations(w http.ResponseWriter, r *http.Request) {
	resp, err := h.dc.GetRecommendationsList(r.Context(), &storagepb.GetRecommendationsListRequest{})
	if err != nil {
		log.Error().Err(err).Msg("get dao recommendations")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	result := make(dao.Recommendations, 0, len(resp.List))
	for _, info := range resp.List {
		result = append(result, dao.Recommendation{
			OriginalId: info.GetOriginalId(),
			InternalId: info.GetInternalId(),
			Name:       info.GetName(),
			Symbol:     info.GetSymbol(),
			NetworkId:  info.GetNetworkId(),
			Address:    info.GetAddress(),
		})
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *DAO) getDelegates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	daoID := vars["id"]

	form, verr := forms.NewGetDelegatesForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetDelegates)

	// TODO: for now we only support one address in query
	var qAccounts []string
	if params.Query != nil {
		qAccounts = append(qAccounts, *params.Query)
	}

	dt := convertDelegationTypeToProto(pointy.StringValue(params.DelegationType, ""))
	resp, err := h.delegateClient.GetDelegates(r.Context(), &storagepb.GetDelegatesRequest{
		DaoId:          daoID,
		QueryAccounts:  qAccounts,
		Sort:           params.By,
		Limit:          int32(params.Limit),
		Offset:         int32(params.Offset),
		DelegationType: &dt,
		ChainId:        params.ChainID,
	})
	if err != nil {
		log.Error().Err(err).Msg("get dao delegates")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	convertedDelegates := make([]dao.Delegate, 0, len(resp.Delegates))
	for _, info := range resp.Delegates {
		convertedDelegates = append(convertedDelegates, dao.Delegate{
			Address:               info.GetAddress(),
			ENSName:               info.GetEnsName(),
			DelegatorCount:        info.GetDelegatorCount(),
			PercentOfDelegators:   info.GetPercentOfDelegators(),
			VotingPower:           info.GetVotingPower(),
			PercentOfVotingPower:  info.GetPercentOfVotingPower(),
			About:                 info.GetAbout(),
			Statement:             info.GetStatement(),
			VotesCount:            info.GetVotesCount(),
			CreatedProposalsCount: info.GetCreatedProposalsCount(),
			DelegationType:        convertDelegationTypeToInternal(info.GetDelegationType()),
			ChainID:               info.ChainId,
		})
	}

	result := dao.DelegatesResponse{
		Delegates: convertedDelegates,
		Total:     resp.Total,
	}

	_ = json.NewEncoder(w).Encode(result)
}

func convertDelegationTypeToProto(value string) storagepb.DelegationType {
	switch value {
	case delegationTypeDelegation:
		return storagepb.DelegationType_DELEGATION_TYPE_DELEGATION
	case delegationTypeErc20Votes:
		return storagepb.DelegationType_DELEGATION_TYPE_ERC20_VOTES
	case delegationTypeSplitDelegation:
		return storagepb.DelegationType_DELEGATION_TYPE_SPLIT_DELEGATION
	default:
		return storagepb.DelegationType_DELEGATION_TYPE_UNRECOGNIZED
	}
}

func convertDelegationTypeToInternal(value storagepb.DelegationType) string {
	switch value {
	case storagepb.DelegationType_DELEGATION_TYPE_DELEGATION:
		return delegationTypeDelegation
	case storagepb.DelegationType_DELEGATION_TYPE_ERC20_VOTES:
		return delegationTypeErc20Votes
	default:
		return delegationTypeSplitDelegation
	}
}

func (h *DAO) getDelegateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	daoID := vars["id"]

	form, verr := forms.NewGetDelegateProfileForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetDelegateProfile)

	resp, err := h.delegateClient.GetDelegateProfile(r.Context(), &storagepb.GetDelegateProfileRequest{
		DaoId:          daoID,
		Address:        params.Address,
		DelegationType: convertDelegationTypeToProto(pointy.StringValue(params.DelegationType, "")),
		ChainId:        params.ChainID,
	})
	if err != nil {
		log.Error().Err(err).Msg("get delegate profile")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	delegates := make([]dao.ProfileDelegateItem, 0, len(resp.Delegates))
	for _, info := range resp.Delegates {
		delegates = append(delegates, dao.ProfileDelegateItem{
			Address:        info.GetAddress(),
			ENSName:        info.GetEnsName(),
			Weight:         info.GetWeight(),
			DelegatedPower: info.GetDelegatedPower(),
		})
	}

	var expiration *time.Time
	if resp.GetExpiration() != nil {
		exp := resp.GetExpiration().AsTime()
		expiration = &exp
	}

	result := dao.DelegateProfile{
		Address:              resp.GetAddress(),
		VotingPower:          resp.GetVotingPower(),
		IncomingPower:        resp.GetIncomingPower(),
		OutgoingPower:        resp.GetOutgoingPower(),
		PercentOfVotingPower: resp.GetPercentOfVotingPower(),
		PercentOfDelegators:  resp.GetPercentOfDelegators(),
		Delegates:            delegates,
		Expiration:           expiration,
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *DAO) getDelegators(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	daoID := vars["id"]
	address := vars["address"]

	form, verr := forms.NewGetDelegatorsForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetDelegators)
	resp, err := h.delegateClient.GetDelegators(r.Context(), &storagepb.GetDelegatorsRequest{
		DaoId:   daoID,
		Address: address,
		ChainId: params.ChainID,
		Limit:   uint32(params.Limit),
		Offset:  helpers.Ptr(uint32(params.Offset)),
	})
	if err != nil {
		log.Error().Err(err).Msg("get delegators")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	convertedDelegators := make([]dao.Delegator, 0, len(resp.GetList()))
	for _, info := range resp.GetList() {
		convertedDelegators = append(convertedDelegators, dao.Delegator{
			Address:    info.GetAddress(),
			ENSName:    info.GetEnsName(),
			TokenValue: info.GetTokenValue(),
		})
	}

	response.AddPaginationHeaders(w, params.Offset, params.Limit, uint64(resp.GetTotalCount()))

	_ = json.NewEncoder(w).Encode(convertedDelegators)
}

func (h *DAO) getTokenInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.dc.GetTokenInfo(r.Context(), &storagepb.TokenInfoRequest{DaoId: id})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"id": id,
		}).Msg("get token info by dao id")
		response.HandleError(response.ResolveError(err), w)
		return
	}

	convertedChains := make([]dao.TokenChainInfo, 0, len(resp.Chains))
	for _, info := range resp.Chains {
		convertedChains = append(convertedChains, dao.TokenChainInfo{
			ChainID:  info.GetChainId(),
			Name:     info.GetName(),
			Decimals: int(info.GetDecimals()),
			IconURL:  info.GetIconUrl(),
			Address:  info.GetAddress(),
		})
	}

	ti := dao.TokenInfo{
		Name:                  resp.GetName(),
		Symbol:                resp.GetSymbol(),
		TotalSupply:           resp.GetTotalSupply(),
		CirculatingSupply:     resp.GetCirculatingSupply(),
		MarketCap:             resp.GetMarketCap(),
		FullyDilutedValuation: resp.GetFullyDilutedValuation(),
		Price:                 resp.GetPrice(),
		FungibleID:            resp.GetFungibleId(),
		Chains:                convertedChains,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ti)
}

func (h *DAO) getTokenChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	period := r.FormValue("period")

	resp, err := h.dc.GetTokenChart(r.Context(), &storagepb.TokenChartRequest{DaoId: id, Period: period})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"id":     id,
			"period": period,
		}).Msg("get token chart by dao id")
		response.HandleError(response.ResolveError(err), w)
		return
	}

	convertedPoints := make([]dao.Point, 0, len(resp.Points))
	for _, info := range resp.Points {
		convertedPoints = append(convertedPoints, dao.Point{
			Time:  info.GetTime().AsTime(),
			Price: info.GetPrice(),
		})
	}

	tc := dao.TokenChart{
		Price:        resp.GetPrice(),
		PriceChanges: resp.GetPriceChanges(),
		Points:       convertedPoints,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tc)
}

func (h *DAO) populateTokenPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.dc.PopulateTokenPrices(r.Context(), &storagepb.TokenPricesRequest{DaoId: id})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"id": id,
		}).Msg("populate token price by dao id")
		response.HandleError(response.ResolveError(err), w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp.Status)
}

func (h *DAO) updateFungibleIds(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue("category")

	resp, err := h.dc.UpdateFungibleIds(r.Context(), &storagepb.UpdateFungibleIdsRequest{Category: category})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"category": category,
		}).Msg("update fungible ids")
		response.HandleError(response.ResolveError(err), w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp.Status)
}

func convertToDaoFromProto(info *storagepb.DaoInfo) dao.Dao {
	id, _ := uuid.Parse(info.GetId())

	return dao.Dao{
		ID:                 id,
		Alias:              info.GetAlias(),
		CreatedAt:          info.GetCreatedAt().AsTime(),
		UpdatedAt:          info.GetUpdatedAt().AsTime(),
		Name:               info.GetName(),
		Private:            info.GetPrivate(),
		About:              info.GetAbout(),
		Avatar:             info.GetAvatar(),
		Terms:              info.GetTerms(),
		Location:           info.GetLocation(),
		Website:            info.GetWebsite(),
		Twitter:            info.GetTwitter(),
		Github:             info.GetGithub(),
		Coingecko:          info.GetCoingeko(),
		Email:              info.GetEmail(),
		Network:            info.GetNetwork(),
		Symbol:             info.GetSymbol(),
		Skin:               info.GetSkin(),
		Domain:             info.GetDomain(),
		Strategies:         convertToStrategiesFromProto(info.GetStrategies()),
		Voting:             convertToVotingFromProto(info.GetVoting()),
		Categories:         info.GetCategories(),
		Treasures:          convertToTreasuresFromProto(info.GetTreasuries()),
		FollowersCount:     info.GetFollowersCount(),
		ProposalsCount:     info.GetProposalsCount(),
		Guidelines:         info.GetGuidelines(),
		Template:           info.GetTemplate(),
		ParentID:           info.GetParentId(),
		ActivitySince:      info.GetActivitySince(),
		VotersCount:        info.GetVotersCount(),
		ActiveVotes:        info.GetActiveVotes(),
		ActiveProposalsIDs: info.GetActiveProposalsIds(),
		Verified:           info.Verified,
		PopularityIndex:    info.GetPopularityIndex(),
		TokenExist:         info.GetTokenExist(),
		TokenSymbol:        info.GetTokenSymbol(),
		FungibleID:         info.GetFungibleId(),
	}
}

func convertToStrategiesFromProto(info []*storagepb.Strategy) dao.Strategies {
	res := make(dao.Strategies, len(info))

	for i, details := range info {
		var params map[string]interface{}
		_ = json.Unmarshal(details.GetParams(), &params)

		res[i] = dao.Strategy{
			Name:    details.GetName(),
			Network: details.GetNetwork(),
			Params:  params,
		}
	}

	return res
}

func convertToTreasuresFromProto(info []*storagepb.Treasury) dao.Treasuries {
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

func convertToVotingFromProto(info *storagepb.Voting) dao.Voting {
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
