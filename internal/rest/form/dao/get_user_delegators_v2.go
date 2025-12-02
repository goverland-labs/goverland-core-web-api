package dao

import (
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
)

type GetUserDelegatorsV2 struct {
	helpers.Pagination

	DelegationType *string
	ChainID        *string
}

func NewGetUserDelegatorsV2Form() *GetUserDelegatorsV2 {
	return &GetUserDelegatorsV2{}
}

func (f *GetUserDelegatorsV2) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	req := GetDelegatesRequest{
		Query:          r.FormValue("query"),
		DelegationType: r.FormValue("delegation_type"),
		ChainID:        r.FormValue("chain_id"),
	}

	f.validateAndSetDelegationType(req, errors)
	f.validateAndSetChainID(req, errors)
	f.ValidateAndSetPagination(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetUserDelegatorsV2) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"offset": f.Offset,
		"limit":  f.Limit,
	}
}

func (f *GetUserDelegatorsV2) validateAndSetDelegationType(req GetDelegatesRequest, _ map[string]response.ErrorMessage) {
	val := strings.TrimSpace(req.DelegationType)
	if val == "" {
		return
	}

	f.DelegationType = &val
}

func (f *GetUserDelegatorsV2) validateAndSetChainID(req GetDelegatesRequest, _ map[string]response.ErrorMessage) {
	value := strings.TrimSpace(req.ChainID)
	if value == "" {
		return
	}

	f.ChainID = &value
}
