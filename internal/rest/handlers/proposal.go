package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	forms "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/proposal"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/proposal"
)

type Proposal struct {
	pc storagepb.ProposalClient
	vc storagepb.VoteClient
}

func NewProposalHandler(pc storagepb.ProposalClient, vc storagepb.VoteClient) APIHandler {
	return &Proposal{
		pc: pc,
		vc: vc,
	}
}

func (h *Proposal) EnrichRoutes(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/proposals/top", h.getTopAction).Methods(http.MethodGet).Name("get_proposals_top")
	baseRouter.HandleFunc("/proposals/{id}/votes", h.getVotesAction).Methods(http.MethodGet).Name("get_proposal_votes")
	baseRouter.HandleFunc("/proposals/{id}/votes/validate", h.validateVote).Methods(http.MethodPost).Name("proposal_vote_validate")
	baseRouter.HandleFunc("/proposals/{id}/votes/prepare", h.prepareVote).Methods(http.MethodPost).Name("proposal_vote_prepare")
	baseRouter.HandleFunc("/proposals/votes", h.vote).Methods(http.MethodPost).Name("proposal_vote")
	baseRouter.HandleFunc("/proposals/{id}", h.getByIDAction).Methods(http.MethodGet).Name("get_proposal_by_id")
	baseRouter.HandleFunc("/proposals", h.getListAction).Methods(http.MethodGet).Name("get_proposals_list")
}

func (h *Proposal) getByIDAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.pc.GetByID(r.Context(), &storagepb.ProposalByIDRequest{ProposalId: id})
	if err != nil {
		log.Error().Err(err).Fields(map[string]interface{}{
			"id": id,
		}).Msg("get proposal by id")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convertToProposalFromProto(resp.Proposal))
}

func (h *Proposal) getListAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewGetListForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetList)
	list, err := h.pc.GetByFilter(r.Context(), &storagepb.ProposalByFilterRequest{
		Dao:         &params.Dao,
		Category:    &params.Category,
		Limit:       &params.Limit,
		Offset:      &params.Offset,
		Title:       &params.Title,
		ProposalIds: params.Proposals,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("get proposal list by filter")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	resp := make([]proposal.Proposal, len(list.GetProposals()))
	for i, info := range list.GetProposals() {
		resp[i] = convertToProposalFromProto(info)
	}

	response.AddPaginationHeaders(w, params.Offset, params.Limit, list.TotalCount)

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Proposal) getTopAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewGetTopForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetTop)
	top := true
	list, err := h.pc.GetByFilter(r.Context(), &storagepb.ProposalByFilterRequest{
		Limit:  &params.Limit,
		Offset: &params.Offset,
		Top:    &top,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("get proposal top by filter")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	resp := make([]proposal.Proposal, len(list.GetProposals()))
	for i, info := range list.GetProposals() {
		resp[i] = convertToProposalFromProto(info)
	}

	response.AddPaginationHeaders(w, params.Offset, params.Limit, list.TotalCount)

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Proposal) getVotesAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	form, verr := forms.NewGetVotesForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetVotes)
	list, err := h.vc.GetVotes(r.Context(), &storagepb.VotesFilterRequest{
		ProposalIds:  []string{id},
		OrderByVoter: &params.Voter,
		Limit:        &params.Limit,
		Offset:       &params.Offset,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("get proposal votes")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	resp := make([]proposal.Vote, len(list.GetVotes()))
	for i, info := range list.GetVotes() {
		resp[i] = convertToProposalVoteFromProto(info)
	}

	response.AddPaginationHeaders(w, params.Offset, params.Limit, list.TotalCount)

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Proposal) validateVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["id"]

	form, verr := forms.NewValidateVoteForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.ValidateVote)
	validateResponse, err := h.vc.Validate(r.Context(), &storagepb.ValidateRequest{
		Voter:    string(params.Voter),
		Proposal: proposalID,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("validate proposal vote")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	var voteValidationError *proposal.VoteValidationError
	if validateResponse.GetValidationError() != nil {
		voteValidationError = &proposal.VoteValidationError{
			Message: validateResponse.GetValidationError().GetMessage(),
			Code:    validateResponse.GetValidationError().GetCode(),
		}
	}

	voteValidation := proposal.VoteValidation{
		OK:                  validateResponse.GetOk(),
		VotingPower:         validateResponse.GetVotingPower(),
		VoteValidationError: voteValidationError,
	}

	_ = json.NewEncoder(w).Encode(voteValidation)
}

func (h *Proposal) prepareVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["id"]

	form, verr := forms.NewPrepareVoteForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.PrepareVote)
	prepareResponse, err := h.vc.Prepare(r.Context(), &storagepb.PrepareRequest{
		Voter:    string(params.Voter),
		Proposal: proposalID,
		Choice: &any.Any{
			Value: params.Choice,
		},
		Reason: params.Reason,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("prepare proposal vote")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	votePreparation := proposal.VotePreparation{
		ID:        prepareResponse.GetId(),
		TypedData: prepareResponse.GetTypedData(),
	}

	_ = json.NewEncoder(w).Encode(votePreparation)
}

func (h *Proposal) vote(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewVoteForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.Vote)
	voteResponse, err := h.vc.Vote(r.Context(), &storagepb.VoteRequest{
		Id:  params.ID,
		Sig: params.Sig,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("vote proposal")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	successfulVote := proposal.SuccessfulVote{
		ID:   voteResponse.GetId(),
		IPFS: voteResponse.GetIpfs(),
		Relayer: proposal.Relayer{
			Address: voteResponse.GetRelayer().GetAddress(),
			Receipt: voteResponse.GetRelayer().GetReceipt(),
		},
	}

	_ = json.NewEncoder(w).Encode(successfulVote)
}

func convertToProposalVoteFromProto(info *storagepb.VoteInfo) proposal.Vote {
	return proposal.Vote{
		ID:           info.GetId(),
		Ipfs:         info.GetIpfs(),
		DaoID:        uuid.MustParse(info.GetDaoId()),
		ProposalID:   info.GetProposalId(),
		Voter:        info.GetVoter(),
		EnsName:      info.GetEnsName(),
		Created:      info.GetCreated(),
		Reason:       info.GetReason(),
		Choice:       info.GetChoice().GetValue(),
		App:          info.GetApp(),
		Vp:           info.GetVp(),
		VpByStrategy: info.GetVpByStrategy(),
		VpState:      info.GetVpState(),
	}
}

func convertToProposalFromProto(info *storagepb.ProposalInfo) proposal.Proposal {
	daoID, _ := uuid.Parse(info.GetDaoId())

	return proposal.Proposal{
		ID:            info.GetId(),
		CreatedAt:     info.GetCreatedAt().AsTime(),
		UpdatedAt:     info.GetUpdatedAt().AsTime(),
		Ipfs:          info.GetIpfs(),
		Author:        info.GetAuthor(),
		EnsName:       info.GetEnsName(),
		Created:       info.GetCreated(),
		DaoID:         daoID,
		Network:       info.GetNetwork(),
		Symbol:        info.GetSymbol(),
		Type:          info.GetType(),
		Strategies:    convertToProposalStrategiesFromProto(info.GetStrategies()),
		Title:         info.GetTitle(),
		Body:          info.GetBody(),
		Discussion:    info.GetDiscussion(),
		Choices:       info.GetChoices(),
		Start:         info.GetStart(),
		End:           info.GetEnd(),
		Quorum:        info.GetQuorum(),
		Privacy:       info.GetPrivacy(),
		Snapshot:      info.GetSnapshot(),
		State:         info.GetState(),
		Link:          info.GetLink(),
		App:           info.GetApp(),
		Scores:        info.GetScores(),
		ScoresState:   info.GetScoresState(),
		ScoresTotal:   info.GetScoresTotal(),
		ScoresUpdated: info.GetScoresUpdated(),
		Votes:         info.GetVotes(),
		Timeline:      convertProposalTimelineFromProto(info.GetTimeline()),
	}
}

func convertToProposalStrategiesFromProto(info []*storagepb.Strategy) proposal.Strategies {
	res := make(proposal.Strategies, len(info))

	for i, details := range info {
		var params map[string]interface{}
		_ = json.Unmarshal(details.GetParams(), &params)

		res[i] = proposal.Strategy{
			Name:    details.GetName(),
			Network: details.GetNetwork(),
			Params:  params,
		}
	}

	return res
}

func convertProposalTimelineFromProto(timeline []*storagepb.ProposalTimelineItem) []proposal.TimelineItem {
	if len(timeline) == 0 {
		return nil
	}

	converted := make([]proposal.TimelineItem, 0, len(timeline))

	for _, t := range timeline {
		converted = append(converted, proposal.TimelineItem{
			CreatedAt: t.GetCreatedAt().AsTime(),
			Action:    convertProposalTimelineActionProto(t.GetAction()),
		})
	}

	return converted
}

var proposalTimelineActionMap = map[storagepb.ProposalTimelineItem_TimelineAction]proposal.TimelineAction{
	storagepb.ProposalTimelineItem_ProposalCreated:             proposal.Created,
	storagepb.ProposalTimelineItem_ProposalUpdated:             proposal.Updated,
	storagepb.ProposalTimelineItem_ProposalVotingStartsSoon:    proposal.VotingStartsSoon,
	storagepb.ProposalTimelineItem_ProposalVotingEndsSoon:      proposal.VotingEndsSoon,
	storagepb.ProposalTimelineItem_ProposalVotingStarted:       proposal.VotingStarted,
	storagepb.ProposalTimelineItem_ProposalVotingQuorumReached: proposal.VotingQuorumReached,
	storagepb.ProposalTimelineItem_ProposalVotingEnded:         proposal.VotingEnded,
}

func convertProposalTimelineActionProto(action storagepb.ProposalTimelineItem_TimelineAction) proposal.TimelineAction {
	converted, exists := proposalTimelineActionMap[action]
	if !exists {
		log.Warn().Str("action", action.String()).Msg("unknown timeline action")
		return proposal.None
	}

	return converted
}
