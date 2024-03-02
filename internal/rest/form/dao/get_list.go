package dao

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
)

type GetListRequest struct {
	Query    string
	Category string
	DAOs     string
}

type GetList struct {
	helpers.Pagination

	Query    *string
	Category *string
	DAOs     []string
}

func NewGetListForm() *GetList {
	return &GetList{}
}

func (f *GetList) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	req := GetListRequest{
		Query:    r.FormValue("query"),
		Category: r.FormValue("category"),
		DAOs:     r.FormValue("daos"),
	}

	f.validateAndSetQuery(req, errors)
	f.validateAndSetCategory(req, errors)
	f.validateAndSetDAOs(req, errors)
	f.ValidateAndSetPagination(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetList) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"query":    f.Query,
		"category": f.Category,
		"offset":   f.Offset,
		"limit":    f.Limit,
	}
}

func (f *GetList) validateAndSetQuery(req GetListRequest, _ map[string]response.ErrorMessage) {
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return
	}

	f.Query = &query
}

func (f *GetList) validateAndSetCategory(req GetListRequest, _ map[string]response.ErrorMessage) {
	category := strings.TrimSpace(req.Category)
	if category == "" {
		return
	}

	f.Category = &category
}

func (f *GetList) validateAndSetDAOs(req GetListRequest, errors map[string]response.ErrorMessage) {
	idsString := strings.TrimSpace(req.DAOs)
	if idsString == "" {
		return
	}

	ids := strings.Split(idsString, ",")
	daos := make([]string, 0, len(ids))
	for i := range ids {
		id := strings.TrimSpace(ids[i])
		if id == "" {
			errors[fmt.Sprintf("daos.%d", i)] = response.WrongValueError("wrong value")
			continue
		}

		// TODO: uuid?
		daos = append(daos, id)
	}

	f.DAOs = daos
}
