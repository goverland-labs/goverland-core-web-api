package proposal

import (
	"net/http"

	"github.com/goverland-labs/core-web-api/internal/response"
	"github.com/goverland-labs/core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/core-web-api/internal/rest/form/common"
)

type GetList struct {
	helpers.Pagination

	Dao      string
	Category string
}

func NewGetListForm() *GetList {
	return &GetList{}
}

func (f *GetList) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	f.Dao = r.FormValue("dao")
	f.Category = r.FormValue("category")
	f.ValidateAndSetPagination(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetList) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"dao":      f.Dao,
		"category": f.Category,
		"offset":   f.Offset,
		"limit":    f.Limit,
	}
}
