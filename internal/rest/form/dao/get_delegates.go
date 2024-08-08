package dao

import (
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
)

type GetDelegatesRequest struct {
	Query string
	By    string
}

type GetDelegates struct {
	helpers.Pagination

	Query *string
	By    *string
}

func NewGetDelegatesForm() *GetDelegates {
	return &GetDelegates{}
}

func (f *GetDelegates) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	req := GetDelegatesRequest{
		Query: r.FormValue("query"),
		By:    r.FormValue("by"),
	}

	f.validateAndSetQuery(req, errors)
	f.validateAndSetBy(req, errors)
	f.ValidateAndSetPagination(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetDelegates) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"query":  f.Query,
		"by":     f.By,
		"offset": f.Offset,
		"limit":  f.Limit,
	}
}

func (f *GetDelegates) validateAndSetQuery(req GetDelegatesRequest, _ map[string]response.ErrorMessage) {
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return
	}

	f.Query = &query
}

func (f *GetDelegates) validateAndSetBy(req GetDelegatesRequest, _ map[string]response.ErrorMessage) {
	by := strings.TrimSpace(req.By)
	if by == "" {
		return
	}

	f.By = &by
}
