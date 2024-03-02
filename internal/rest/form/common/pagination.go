package common

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/response/errs"
)

const (
	DefaultOffset uint64 = 0
	DefaultLimit  uint64 = 20
)

type Pagination struct {
	Offset uint64
	Limit  uint64
}

func (p *Pagination) ValidateAndSetPagination(r *http.Request, errors map[string]response.ErrorMessage) {
	p.validateAndSetOffset(r, errors)
	p.validateAndSetLimit(r, errors)
}

func (p *Pagination) ValidateAndSetPaginationWithLimit(
	r *http.Request,
	maxLimit uint64,
	errors map[string]response.ErrorMessage,
) {
	p.validateAndSetOffset(r, errors)
	p.validateAndSetLimitMax(r, maxLimit, errors)
}

func (p *Pagination) validateAndSetOffset(r *http.Request, errors map[string]response.ErrorMessage) {
	offset := r.FormValue("offset")
	if offset == "" {
		p.Offset = DefaultOffset

		return
	}

	number, err := strconv.ParseInt(offset, 10, 64) // nolint:gomnd
	if err != nil {
		errors["offset"] = response.ErrorMessage{
			Code:    errs.WrongFormat,
			Message: "should be integer",
		}

		return
	}

	if number < 0 {
		errors["offset"] = response.ErrorMessage{
			Code:    errs.WrongValue,
			Message: "should be more than 0",
		}

		return
	}

	p.Offset = uint64(number)
}

func (p *Pagination) validateAndSetLimit(r *http.Request, errors map[string]response.ErrorMessage) {
	limit := r.FormValue("limit")
	if limit == "" {
		p.Limit = DefaultLimit

		return
	}

	number, err := strconv.ParseInt(limit, 10, 64) // nolint:gomnd
	if err != nil {
		errors["limit"] = response.ErrorMessage{
			Code:    errs.WrongFormat,
			Message: "should be integer",
		}

		return
	}

	if number <= 0 {
		errors["limit"] = response.ErrorMessage{
			Code:    errs.WrongValue,
			Message: "should be more than 0",
		}

		return
	}

	p.Limit = uint64(number)
}

func (p *Pagination) validateAndSetLimitMax(r *http.Request, maxLimit uint64, errors map[string]response.ErrorMessage) {
	limit := r.FormValue("limit")
	if limit == "" {
		p.Limit = p.minUint64(DefaultLimit, maxLimit)

		return
	}

	number, err := strconv.ParseInt(limit, 10, 64) // nolint:gomnd
	if err != nil {
		errors["limit"] = response.ErrorMessage{
			Code:    errs.WrongFormat,
			Message: "should be integer",
		}

		return
	}

	if number <= 0 {
		errors["limit"] = response.ErrorMessage{
			Code:    errs.WrongValue,
			Message: "should be more than 0",
		}

		return
	}

	if uint64(number) > maxLimit {
		errors["limit"] = response.ErrorMessage{
			Code:    errs.WrongValue,
			Message: fmt.Sprintf("should be less or equal than %d", maxLimit),
		}

		return
	}

	p.Limit = uint64(number)
}

func (p *Pagination) minUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}

	return b
}
