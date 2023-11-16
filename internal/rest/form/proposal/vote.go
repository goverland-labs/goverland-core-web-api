package proposal

import (
	"encoding/json"
	"net/http"

	"github.com/goverland-labs/core-web-api/internal/response"
	"github.com/goverland-labs/core-web-api/internal/response/errs"
	"github.com/goverland-labs/core-web-api/internal/rest/form"
)

type VoteRequest struct {
	ID  uint64 `json:"id"`
	Sig string `json:"sig"`
}

type Vote struct {
	ID  uint64 `json:"id"`
	Sig string `json:"sig"`
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

	f.ID = req.ID
	f.Sig = req.Sig

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *Vote) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":  f.ID,
		"sig": f.Sig,
	}
}
