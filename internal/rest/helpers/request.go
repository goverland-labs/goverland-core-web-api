package helpers

import (
	"net/http"
	"strings"
)

const (
	tokenHeader = "Authorization"
	tokenPrefix = ""
)

func ExtractTokenFromHeaders(header http.Header) string {
	value := header.Get(tokenHeader)

	if value == "" {
		return ""
	}

	if !strings.HasPrefix(value, tokenPrefix) {
		return ""
	}

	return strings.TrimPrefix(value, tokenPrefix)
}
