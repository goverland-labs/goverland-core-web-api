package response

import (
	"math"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/core-web-api/internal/response/errs"
)

type parametrizedError interface {
	SetError(key string, code errs.ErrCode, message string)
}

func ResolveError(err error) Error {
	details, ok := status.FromError(err)
	if !ok {
		return NewInternalError()
	}

	switch details.Code() {
	case codes.InvalidArgument:
		ed := NewValidationError()
		enrichErrDetails(ed, details)
		return ed
	case codes.NotFound:
		return NewNotFoundError()

	case codes.PermissionDenied:
		return NewPermissionDeniedError()

	case codes.ResourceExhausted:
		retryAfter := 0
		for _, d := range details.Details() {
			if info, ok := d.(*errdetails.RetryInfo); ok {
				delay := info.GetRetryDelay().AsDuration()
				retryAfter = int(math.Ceil(delay.Seconds()))
			}
		}
		return NewRateLimitedError(retryAfter)
	}

	return NewInternalError()
}

func enrichErrDetails(err parametrizedError, st *status.Status) {
	err.SetError(GeneralErrorKey, errs.WrongValue, st.Message())

	return
}

// IsInternalError returns false when the error is caused by invalid request data or by other mismatches caused by an user.
// It can be used to determine whether the error should be logged. If the function returns false, the error shouldn't
// be logged.
func IsInternalError(err error) bool {
	if err == nil {
		return false
	}

	e := ResolveError(err)
	switch e.(type) {
	case *ValidationError:
		return false

	case *RateLimitedError:
		return false
	}

	return true
}
