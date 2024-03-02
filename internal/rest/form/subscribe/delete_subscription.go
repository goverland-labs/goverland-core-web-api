package subscribe

import (
	"encoding/json"
	"net/http"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/response/errs"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
)

type UnsubscribeOnDaoRequest struct {
	DaoID string `json:"dao"`
}

type UnsubscribeOnDaoForm struct {
	DaoID string
}

func NewUnsubscribeOnDaoForm() *UnsubscribeOnDaoForm {
	return &UnsubscribeOnDaoForm{}
}

func (f *UnsubscribeOnDaoForm) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	var req *UnsubscribeOnDaoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, errs.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	f.DaoID = req.DaoID

	return f, nil
}

func (f *UnsubscribeOnDaoForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"dao": f.DaoID,
	}
}
