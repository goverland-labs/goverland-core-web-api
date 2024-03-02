package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/goverland-labs/goverland-core-feed/protocol/feedpb"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	forms "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/feed"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/dao"
)

type Feed struct {
	fc feedpb.FeedClient
}

func NewFeedHandler(fc feedpb.FeedClient) APIHandler {
	return &Feed{
		fc: fc,
	}
}

func (h *Feed) EnrichRoutes(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/feed", h.getFeedByFiltersAction).Methods(http.MethodPost).Name("get_feed_by_filters")
}

func (h *Feed) getFeedByFiltersAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewGetFeedListForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.GetFeedList)
	resp, err := h.fc.GetByFilter(r.Context(), &feedpb.FeedByFilterRequest{
		DaoIds:   params.DaoList,
		Types:    params.Types,
		Actions:  params.Actions,
		IsActive: params.IsActive,
		Limit:    &params.Limit,
		Offset:   &params.Offset,
	})
	if err != nil {
		log.Error().Err(err).Msg("get feed by filters")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	list := make([]dao.FeedItem, len(resp.Items))
	for i, fi := range resp.Items {
		list[i] = convertToFeedItemFromProto(fi)
	}

	response.AddPaginationHeaders(w, params.Offset, params.Limit, resp.TotalCount)
	_ = json.NewEncoder(w).Encode(list)
}
