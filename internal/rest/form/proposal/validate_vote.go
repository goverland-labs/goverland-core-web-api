package proposal

import (
	"encoding/json"
	"net/http"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/response/errs"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
)

type ValidateVoteRequest struct {
	Voter string `json:"voter"`
}

type ValidateVote struct {
	Voter common.Voter `json:"voter"`
}

func NewValidateVoteForm() *ValidateVote {
	return &ValidateVote{}
}

func (f *ValidateVote) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	var req *ValidateVoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, errs.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errors := make(map[string]response.ErrorMessage)

	f.Voter.ValidateAndSet(req.Voter, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *ValidateVote) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"voter": f.Voter,
	}
}
