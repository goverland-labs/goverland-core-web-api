package common

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
)

type GetUserVotes struct {
	Pagination

	Proposals []string
}

func NewGetUserVotesForm() *GetUserVotes {
	return &GetUserVotes{}
}

func (f *GetUserVotes) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	f.ValidateAndSetPagination(r, errors)
	f.validateAndSetProposalIds(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetUserVotes) validateAndSetProposalIds(r *http.Request, errors map[string]response.ErrorMessage) {
	pr := strings.TrimSpace(r.FormValue("proposals"))
	if pr == "" {
		return
	}

	ids := strings.Split(pr, ",")
	proposals := make([]string, 0, len(ids))
	for i := range ids {
		id := strings.TrimSpace(ids[i])
		if id == "" {
			errors[fmt.Sprintf("proposals.%d", i)] = response.WrongValueError("wrong value")
			continue
		}

		proposals = append(proposals, id)
	}

	f.Proposals = proposals
}

func (f *GetUserVotes) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"offset": f.Offset,
		"limit":  f.Limit,
	}
}
