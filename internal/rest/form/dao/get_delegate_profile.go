package dao

import (
	"net/http"
	"strings"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/form"
)

type GetDelegateProfileRequest struct {
	Address        string
	DelegationType string
	ChainID        string
}

type GetDelegateProfile struct {
	Address        string
	DelegationType *string
	ChainID        *string
}

func NewGetDelegateProfileForm() *GetDelegateProfile {
	return &GetDelegateProfile{}
}

func (f *GetDelegateProfile) ParseAndValidate(r *http.Request) (form.Former, response.Error) {
	errors := make(map[string]response.ErrorMessage)

	req := GetDelegateProfileRequest{
		Address:        r.FormValue("address"),
		DelegationType: r.FormValue("delegation_type"),
		ChainID:        r.FormValue("chain_id"),
	}

	f.validateAndSetAddress(req, errors)
	f.validateAndSetDelegationType(req, errors)
	f.validateAndSetChainID(req, errors)

	if len(errors) > 0 {
		ve := response.NewValidationError(errors)

		return nil, ve
	}

	return f, nil
}

func (f *GetDelegateProfile) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"address":         f.Address,
		"delegation_type": f.DelegationType,
		"chain_id":        f.ChainID,
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

func (f *GetDelegateProfile) validateAndSetDelegationType(req GetDelegateProfileRequest, _ map[string]response.ErrorMessage) {
	val := strings.TrimSpace(req.DelegationType)
	if val == "" {
		return
	}

	f.DelegationType = &val
}

func (f *GetDelegateProfile) validateAndSetChainID(req GetDelegateProfileRequest, _ map[string]response.ErrorMessage) {
	value := strings.TrimSpace(req.ChainID)
	if value == "" {
		return
	}

	f.ChainID = &value
}
