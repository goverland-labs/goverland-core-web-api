package subscribe

import (
	"encoding/json"
	"net/http"

	"github.com/goverland-labs/core-web-api/internal/response"
	"github.com/goverland-labs/core-web-api/internal/response/errs"
	"github.com/goverland-labs/core-web-api/internal/rest/form"
)

type SubscribeOnDaoRequest struct {
	DaoID string `json:"dao"`
}

type SubscribeOnDaoForm struct {
	DaoID string
}

func NewSubscribeOnDaoForm() *SubscribeOnDaoForm {
	return &SubscribeOnDaoForm{}
}

func (f *SubscribeOnDaoForm) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	var req *SubscribeOnDaoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, errs.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	f.DaoID = req.DaoID

	return f, nil
}

func (f *SubscribeOnDaoForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"dao": f.DaoID,
	}
}
