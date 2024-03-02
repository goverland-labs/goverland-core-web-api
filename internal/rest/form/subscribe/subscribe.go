package subscribe

import (
	"encoding/json"
	"net/http"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/response/errs"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
)

type SubscribeRequest struct {
	WebhookURL string `json:"webhook_url"`
}

type SubscribeForm struct {
	WebhookURL string
}

func NewSubscribeForm() *SubscribeForm {
	return &SubscribeForm{}
}

func (f *SubscribeForm) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	var req *SubscribeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ve := response.NewValidationError()
		ve.SetError(response.GeneralErrorKey, errs.InvalidRequestStructure, "invalid request structure")

		return nil, ve
	}

	f.WebhookURL = req.WebhookURL

	return f, nil
}

func (f *SubscribeForm) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"webhook_url": f.WebhookURL,
	}
}
