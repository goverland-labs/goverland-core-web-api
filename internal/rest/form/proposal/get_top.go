package proposal

import (
	"net/http"

	"github.com/goverland-labs/core-web-api/internal/response"
	"github.com/goverland-labs/core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/core-web-api/internal/rest/form/common"
)

type GetTop struct {
	helpers.Pagination
}

func NewGetTopForm() *GetTop {
	return &GetTop{}
}

func (f *GetTop) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	f.ValidateAndSetPagination(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetTop) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"offset": f.Offset,
		"limit":  f.Limit,
	}
}
