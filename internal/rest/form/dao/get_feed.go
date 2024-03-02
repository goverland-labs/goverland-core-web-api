package dao

import (
	"net/http"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
)

type GetFeed struct {
	helpers.Pagination
}

func NewGetFeedForm() *GetFeed {
	return &GetFeed{}
}

func (f *GetFeed) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	f.ValidateAndSetPagination(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetFeed) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"offset": f.Offset,
		"limit":  f.Limit,
	}
}
