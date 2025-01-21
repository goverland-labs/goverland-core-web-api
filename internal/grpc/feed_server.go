package grpc

import (
	"fmt"
	"time"

	"github.com/goverland-labs/goverland-core-web-api/protocol/feed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FeedServer struct {
	feed.UnimplementedFeedEventsServer

	service *Service
}

func NewFeedServer(sp *Service) *FeedServer {
	return &FeedServer{
		service: sp,
	}
}

func (s *FeedServer) EventsSubscribe(req *feed.EventsSubscribeRequest, stream grpc.ServerStreamingServer[feed.FeedItem]) error {
	ctx := stream.Context()

	var lastUpdated *time.Time
	if req.GetLastUpdatedAt() != nil {
		lu := req.GetLastUpdatedAt().AsTime()
		lastUpdated = &lu
	}

	events := s.service.GetFeedItems(ctx, ItemsRequest{
		SubscriberID:      req.GetSubscriberId(),
		SubscriptionTypes: req.GetSubscriptionTypes(),
		LastUpdatedAt:     lastUpdated,
	})

	for {
		var (
			event Result
			ok    bool
		)

		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		case event, ok = <-events:
			if !ok {
				return status.Error(codes.Canceled, "events channel closed")
			}

			if event.Err != nil {
				return status.Errorf(codes.Internal, "internal error on getting data: %v", event.Err)
			}
		}

		if err := stream.Send(event.Item); err != nil {
			return fmt.Errorf("send event: %w", err)
		}
	}
}
