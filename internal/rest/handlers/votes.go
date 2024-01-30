package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/goverland-labs/core-api/protobuf/internalapi"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/core-web-api/internal/response"
	forms "github.com/goverland-labs/core-web-api/internal/rest/form/common"
	"github.com/goverland-labs/core-web-api/internal/rest/models/proposal"
)

type Votes struct {
	vc internalapi.VoteClient
}

func NewVotesHandler(vc internalapi.VoteClient) APIHandler {
	return &Votes{
		vc: vc,
	}
}

func (h *Votes) EnrichRoutes(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/user/{address}/votes", h.getUserVotesAction).Methods(http.MethodGet).Name("get_user_votes")
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
	var list, err = h.vc.GetVotes(r.Context(), &internalapi.VotesFilterRequest{
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

	_ = json.NewEncoder(w).Encode(resp)
}

func convertToVoteFromProto(info *internalapi.VoteInfo) proposal.Vote {
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
