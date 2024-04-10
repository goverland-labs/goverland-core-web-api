package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	forms "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/proposal"
)

type Votes struct {
	vc storagepb.VoteClient
}

func NewVotesHandler(vc storagepb.VoteClient) APIHandler {
	return &Votes{
		vc: vc,
	}
}

func (h *Votes) EnrichRoutes(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/user/{address}/votes", h.getUserVotesAction).Methods(http.MethodGet).Name("get_user_votes")
	baseRouter.HandleFunc("/user/{address}/participated-daos", h.getUserParticipatedDaos).Methods(http.MethodGet).Name("get_user_participated_daos")
}

func (h *Votes) getUserVotesAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	form, verr := forms.NewGetUserVotesForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetUserVotes)
	var list, err = h.vc.GetVotes(r.Context(), &storagepb.VotesFilterRequest{
		ProposalIds: params.Proposals,
		Voter:       &address,
		Limit:       &params.Limit,
		Offset:      &params.Offset,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("get user votes")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	resp := make([]proposal.Vote, len(list.GetVotes()))
	for i, info := range list.GetVotes() {
		resp[i] = convertToVoteFromProto(info)
	}

	response.AddPaginationHeaders(w, params.Offset, params.Limit, list.TotalCount)
	response.AddTotalVpHeader(w, list.TotalVp)

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Votes) getUserParticipatedDaos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	var list, err = h.vc.GetDaosVotedIn(r.Context(), &storagepb.DaosVotedInRequest{
		Voter: address,
	})
	if err != nil {
		log.Error().Err(err).Msg("get user participated daos")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	resp := make([]uuid.UUID, len(list.GetDaoIds()))
	for i, info := range list.GetDaoIds() {
		id, _ := uuid.Parse(info)
		resp[i] = id
	}

	response.AddPaginationHeaders(w, 0, list.TotalCount, list.TotalCount)
	_ = json.NewEncoder(w).Encode(resp)
}

func convertToVoteFromProto(info *storagepb.VoteInfo) proposal.Vote {
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
