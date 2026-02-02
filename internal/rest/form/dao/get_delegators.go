package dao

import (
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
)

type GetDelegatorsRequest struct {
	ChainID        string
	DelegationType string
}

type GetDelegators struct {
	ChainID        string
	DelegationType string

	helpers.Pagination
}

func NewGetDelegatorsForm() *GetDelegators {
	return &GetDelegators{}
}

func (f *GetDelegators) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)
	req := GetDelegatorsRequest{
		ChainID:        r.FormValue("chain_id"),
		DelegationType: r.FormValue("delegation_type"),
	}

	f.validateAndSetChainID(req, errors)
	f.validateAndSetDelegationType(req, errors)
	f.ValidateAndSetPagination(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)
		return nil, ve
	}

	return f, nil
}

func (f *GetDelegators) validateAndSetDelegationType(req GetDelegatorsRequest, _ map[string]response.ErrorMessage) {
	value := strings.TrimSpace(req.DelegationType)
	if value == "" {
		return
	}

	f.DelegationType = value
}

func (f *GetDelegators) validateAndSetChainID(req GetDelegatorsRequest, _ map[string]response.ErrorMessage) {
	value := strings.TrimSpace(req.ChainID)
	if value == "" {
		return
	}

	f.ChainID = value
}

func (f *GetDelegators) ConvertToMap() map[string]interface{} {
	return map[string]any{
		"chain_id":        f.ChainID,
		"delegation_type": f.DelegationType,
		"offset":          f.Offset,
		"limit":           f.Limit,
	}
}
