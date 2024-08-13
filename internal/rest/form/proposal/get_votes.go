package proposal

import (
	"net/http"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
	helpers "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/common"
)

type GetVotes struct {
	helpers.Pagination

	Voter string
	Query string
}

func NewGetVotesForm() *GetVotes {
	return &GetVotes{}
}

func (f *GetVotes) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	f.ValidateAndSetPagination(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}
	f.Voter = r.FormValue("voter")
	f.Query = r.FormValue("query")

	return f, nil
}

func (f *GetVotes) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"voter":  f.Voter,
		"query":  f.Query,
		"offset": f.Offset,
		"limit":  f.Limit,
	}
}
