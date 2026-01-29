package ens

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
)

type GetAddressesByNames struct {
	Names []string
}

func NewGetAddressesByNamesForm() *GetAddressesByNames {
	return &GetAddressesByNames{}
}

func (f *GetAddressesByNames) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	f.validateAndSetNames(r, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetAddressesByNames) validateAndSetNames(r *http.Request, errors map[string]response.ErrorMessage) {
	raw := strings.TrimSpace(r.FormValue("names"))
	if raw == "" {
		return
	}

	params := strings.Split(raw, ",")
	names := make([]string, 0, len(params))
	for i := range params {
		name := strings.TrimSpace(params[i])
		if name == "" {
			errors[fmt.Sprintf("names.%d", i)] = response.WrongValueError("wrong value")
			continue
		}

		names = append(names, name)
	}

	f.Names = names
}

func (f *GetAddressesByNames) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{}
}
