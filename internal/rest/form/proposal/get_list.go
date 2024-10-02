package proposal

import (
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
)

type GetList struct {
	helpers.Pagination

	Dao        string
	Category   string
	Title      string
	Proposals  []string
	OnlyActive bool
}

func NewGetListForm() *GetList {
	return &GetList{}
}

func (f *GetList) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	f.Dao = r.FormValue("dao")
	f.Category = r.FormValue("category")
	f.Title = r.FormValue("title")
	f.OnlyActive = r.FormValue("only_active") == "true"
	idsString := strings.TrimSpace(r.FormValue("proposals"))
	if idsString != "" {
		ids := strings.Split(idsString, ",")
		f.Proposals = ids
	}
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
		"title":    f.Title,
		"offset":   f.Offset,
		"limit":    f.Limit,
	}
}
