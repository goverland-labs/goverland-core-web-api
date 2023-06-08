package response

import (
	"encoding/json"
	"net/http"
)

func HandleError(err Error, w http.ResponseWriter) {
	if err == nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(err.GetHTTPStatus())
	_ = json.NewEncoder(w).Encode(ParseError(err))
}

// ParseError determines the error type and creates a map with the error description.
func ParseError(err Error) map[string]interface{} {
	if err == nil {
		return nil
	}

	switch e := err.(type) { // nolint:gocritic
	case *ValidationError:
		return map[string]interface{}{
			"errors":  e.Errors(),
			"message": e.PublicMessage(),
		}

	case *UnprocessableEntityError:
		return map[string]interface{}{
			"errors":  e.Errors(),
			"message": e.PublicMessage(),
		}

	case *RateLimitedError:
		return map[string]interface{}{
			"message":     e.PublicMessage(),
			"retry_after": e.RetryAfter,
		}

	default:
		return map[string]interface{}{
			"message": err.PublicMessage(),
		}
	}
}
