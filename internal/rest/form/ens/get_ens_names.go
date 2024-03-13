package ens

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
)

type GetEnsNames struct {
	Addresses []string
}

func NewGetEnsNamesForm() *GetEnsNames {
	return &GetEnsNames{}
}

func (f *GetEnsNames) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	f.validateAndSetAddresses(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetEnsNames) validateAndSetAddresses(r *http.Request, errors map[string]response.ErrorMessage) {
	adr := strings.TrimSpace(r.FormValue("addresses"))
	if adr == "" {
		return
	}

	params := strings.Split(adr, ",")
	addresses := make([]string, 0, len(params))
	for i := range params {
		id := strings.TrimSpace(params[i])
		if id == "" {
			errors[fmt.Sprintf("addresses.%d", i)] = response.WrongValueError("wrong value")
			continue
		}

		addresses = append(addresses, id)
	}

	f.Addresses = addresses
}

func (f *GetEnsNames) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{}
}
