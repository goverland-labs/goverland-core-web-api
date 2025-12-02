package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/goverland-labs/goverland-core-feed/protocol/feedpb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"
	forms "github.com/goverland-labs/goverland-core-web-api/internal/rest/form/subscribe"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/helpers"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest/models/subscribe"
)

type Subscriber struct {
	subscribers   feedpb.SubscriberClient
	subscriptions feedpb.SubscriptionClient
}

func NewSubscribeHandler(subscribers feedpb.SubscriberClient, subscriptions feedpb.SubscriptionClient) *Subscriber {
	return &Subscriber{
		subscribers:   subscribers,
		subscriptions: subscriptions,
	}
}

func (h *Subscriber) EnrichRoutes(v1, _ *mux.Router) {
	v1.HandleFunc("/subscribe", h.createSubscriberAction).Methods(http.MethodPost).Name("create_subscriber")
	v1.HandleFunc("/subscribe", h.updateSubscriberAction).Methods(http.MethodPut).Name("update_subscriber")
	v1.HandleFunc("/subscriptions", h.subscribeOnDaoAction).Methods(http.MethodPost).Name("create_subscription")
	v1.HandleFunc("/subscriptions", h.unsubscribeOnDaoAction).Methods(http.MethodDelete).Name("delete_subscription")
}

func (h *Subscriber) createSubscriberAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewSubscribeForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.SubscribeForm)
	resp, err := h.subscribers.Create(r.Context(), &feedpb.CreateSubscriberRequest{WebhookUrl: params.WebhookURL})
	if err != nil {
		log.Error().Err(err).Msg("create subscriber")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(convertToSubscriberFromProto(resp.GetSubscriberId()))
}

// todo: think about location of this function
func prepareOutgoingContext(ctx context.Context, header http.Header) context.Context {
	subscriberID := helpers.ExtractTokenFromHeaders(header)

	// todo: move key to const
	md := metadata.New(map[string]string{"subscriber_id": subscriberID})
	return metadata.NewOutgoingContext(ctx, md)
}

func (h *Subscriber) updateSubscriberAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewSubscribeForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.SubscribeForm)
	_, err := h.subscribers.Update(prepareOutgoingContext(r.Context(), r.Header), &feedpb.UpdateSubscriberRequest{WebhookUrl: params.WebhookURL})
	if err != nil {
		log.Error().Err(err).Msg("update subscriber")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Subscriber) subscribeOnDaoAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewSubscribeOnDaoForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.SubscribeOnDaoForm)
	_, err := h.subscriptions.Subscribe(prepareOutgoingContext(r.Context(), r.Header), &feedpb.SubscribeRequest{
		DaoId: params.DaoID,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("create subscription")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Subscriber) unsubscribeOnDaoAction(w http.ResponseWriter, r *http.Request) {
	form, verr := forms.NewUnsubscribeOnDaoForm().ParseAndValidate(r)
	if verr != nil {
		response.HandleError(verr, w)

		return
	}

	params := form.(*forms.UnsubscribeOnDaoForm)
	_, err := h.subscriptions.Unsubscribe(prepareOutgoingContext(r.Context(), r.Header), &feedpb.UnsubscribeRequest{
		DaoId: params.DaoID,
	})
	if err != nil {
		log.Error().Err(err).Fields(params.ConvertToMap()).Msg("create subscription")

		response.HandleError(response.ResolveError(err), w)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func convertToSubscriberFromProto(subID string) subscribe.Subscriber {
	return subscribe.Subscriber{
		SubscriberID: subID,
	}
}
