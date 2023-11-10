package proposal

import (
	"encoding/json"
	"net/http"

	"github.com/goverland-labs/core-web-api/internal/response"
	"github.com/goverland-labs/core-web-api/internal/response/errs"
	"github.com/goverland-labs/core-web-api/internal/rest/form"
	"github.com/goverland-labs/core-web-api/internal/rest/form/common"
)

type VoteRequest struct {
	Voter  string          `json:"voter"`
	Choice json.RawMessage `json:"choice"`
	Reason *string         `json:"reason,omitempty"`
	Sig    string          `json:"sig"`
}

type Vote struct {
	Voter  common.Voter  `json:"voter"`
	Choice common.Choice `json:"choice"`
	Reason *string       `json:"reason,omitempty"`
	Sig    string        `json:"sig"`
}

func NewVoteForm() *Vote {
	return &Vote{}
}

func (f *Vote) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	var req *VoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, errs.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errors := make(map[string]response.ErrorMessage)

	f.Voter.ValidateAndSet(req.Voter, errors)
	f.Choice.ValidateAndSet(req.Choice, errors)
	f.Reason = req.Reason
	f.Sig = req.Sig

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *Vote) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"voter":  f.Voter,
		"choice": f.Choice,
		"reason": f.Reason,
		"sig":    f.Sig,
	}
}
