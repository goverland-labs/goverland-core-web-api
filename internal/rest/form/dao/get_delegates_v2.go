package dao

import (
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
)

type GetDelegatesV2Request struct {
	Query          string
	By             string
	DelegationType string
	ChainID        string
}

type GetDelegatesV2 struct {
	helpers.Pagination

	Query          *string
	By             *string
	DelegationType *string
	ChainID        *string
}

func NewGetDelegatesV2Form() *GetDelegatesV2 {
	return &GetDelegatesV2{}
}

func (f *GetDelegatesV2) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	req := GetDelegatesRequest{
		Query:          r.FormValue("query"),
		By:             r.FormValue("by"),
		DelegationType: r.FormValue("delegation_type"),
		ChainID:        r.FormValue("chain_id"),
	}

	f.validateAndSetQuery(req, errors)
	f.validateAndSetBy(req, errors)
	f.validateAndSetDelegationType(req, errors)
	f.validateAndSetChainID(req, errors)
	f.ValidateAndSetPagination(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetDelegatesV2) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"query":  f.Query,
		"by":     f.By,
		"offset": f.Offset,
		"limit":  f.Limit,
	}
}

func (f *GetDelegatesV2) validateAndSetQuery(req GetDelegatesRequest, _ map[string]response.ErrorMessage) {
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return
	}

	f.Query = &query
}

func (f *GetDelegatesV2) validateAndSetBy(req GetDelegatesRequest, _ map[string]response.ErrorMessage) {
	by := strings.TrimSpace(req.By)
	if by == "" {
		return
	}

	f.By = &by
}

func (f *GetDelegatesV2) validateAndSetDelegationType(req GetDelegatesRequest, _ map[string]response.ErrorMessage) {
	val := strings.TrimSpace(req.DelegationType)
	if val == "" {
		return
	}

	f.DelegationType = &val
}

func (f *GetDelegatesV2) validateAndSetChainID(req GetDelegatesRequest, _ map[string]response.ErrorMessage) {
	value := strings.TrimSpace(req.ChainID)
	if value == "" {
		return
	}

	f.ChainID = &value
}
