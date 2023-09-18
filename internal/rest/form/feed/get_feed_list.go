package feed

import (
	"encoding/json"
	"net/http"

	"github.com/goverland-labs/core-web-api/internal/response"
	"github.com/goverland-labs/core-web-api/internal/response/errs"
	"github.com/goverland-labs/core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/core-web-api/internal/rest/form/common"
)

type GetFeedRequest struct {
	DaoList  []string `json:"dao_list"`
	IsActive *bool    `json:"is_active,omitempty"`
	Types    []string `json:"types"`
	Actions  []string `json:"actions"`

	helpers.Pagination
}

type GetFeedList struct {
	DaoList  []string
	IsActive *bool
	Types    []string
	Actions  []string

	helpers.Pagination
}

func NewGetFeedListForm() *GetFeedList {
	return &GetFeedList{}
}

func (f *GetFeedList) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	var req *GetFeedRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, errs.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	errors := make(map[string]response.ErrorMessage)

	f.ValidateAndSetPagination(r, errors)
	f.Types = req.Types
	f.Actions = req.Actions
	f.DaoList = req.DaoList
	f.IsActive = req.IsActive

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetFeedList) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"offset": f.Offset,
		"limit":  f.Limit,
	}
}
