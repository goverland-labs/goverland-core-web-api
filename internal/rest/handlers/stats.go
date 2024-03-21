package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/stats"
)

type Stats struct {
	sc storagepb.StatsClient
}

func NewStatsHandler(sc storagepb.StatsClient) APIHandler {
	return &Stats{
		sc: sc,
	}
}

func (h *Stats) EnrichRoutes(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/stats/totals", h.getTotals).Methods(http.MethodGet).Name("get_stats_totals")
}

func (h *Stats) getTotals(w http.ResponseWriter, r *http.Request) {
	var totals, err = h.sc.GetTotals(r.Context(), &storagepb.GetTotalsRequest{})
	if err != nil {
		log.Error().Err(err).Msg("get stats totals")
		response.HandleError(response.ResolveError(err), w)

		return
	}

	_ = json.NewEncoder(w).Encode(convertToTotalFromProto(totals))
}

func convertToTotalFromProto(info *storagepb.GetTotalsResponse) stats.Totals {
	return stats.Totals{
		Dao: stats.Dao{
			Total:         info.GetDao().GetTotal(),
			TotalVerified: info.GetDao().GetTotalVerified(),
		},
		Proposals: stats.Proposals{
			Total: info.GetProposals().GetTotal(),
		},
	}
}
