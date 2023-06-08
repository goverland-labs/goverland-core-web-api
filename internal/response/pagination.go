package response

import (
	"fmt"
	"net/http"
)

const (
	HeaderTotalCount    = "X-Total-Count"
	HeaderCurrentOffset = "X-Offset"
	HeaderLimit         = "X-Limit"
)

func AddPaginationHeaders(w http.ResponseWriter, offset, limit, totalCnt uint64) {
	w.Header().Set(HeaderTotalCount, fmt.Sprintf("%d", totalCnt))
	w.Header().Set(HeaderCurrentOffset, fmt.Sprintf("%d", offset))
	w.Header().Set(HeaderLimit, fmt.Sprintf("%d", limit))
}
