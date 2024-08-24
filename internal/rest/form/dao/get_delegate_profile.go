package dao

import (
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
)

type GetDelegateProfileRequest struct {
	Address string
}

type GetDelegateProfile struct {
	Address string
}

func NewGetDelegateProfileForm() *GetDelegateProfile {
	return &GetDelegateProfile{}
}

func (f *GetDelegateProfile) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	req := GetDelegateProfileRequest{
		Address: r.FormValue("address"),
	}

	f.validateAndSetAddress(req, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetDelegateProfile) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"address": f.Address,
	}
}

func (f *GetDelegateProfile) validateAndSetAddress(req GetDelegateProfileRequest, errors map[string]response.ErrorMessage) {
	address := strings.TrimSpace(req.Address)
	if address == "" {
		errors["address"] = response.MissedValueError("address is required")
		return
	}

	f.Address = address
}
