package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	feedproto "github.com/goverland-labs/goverland-core-feed/protocol/feedpb"
	coreproto "github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	internalproto "github.com/goverland-labs/goverland-core-web-api/protocol/feed"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	coreFeed    feedproto.FeedEventsClient
	coreStorage coreproto.VoteClient
}

type ItemsRequest struct {
	SubscriberID      string
	SubscriptionTypes []internalproto.FeedItemType
	LastUpdatedAt     *time.Time
}

type FeedItem struct{}

type Result struct {
	Item *internalproto.FeedItem
	Err  error
}

func NewService(fc feedproto.FeedEventsClient, cv coreproto.VoteClient) *Service {
	return &Service{
		coreFeed:    fc,
		coreStorage: cv,
	}
}

func (s *Service) GetFeedItems(ctx context.Context, req ItemsRequest) <-chan Result {
	ch := make(chan Result, 100)
	ctx, cancel := context.WithCancel(ctx)

	var updatedAt *timestamppb.Timestamp
	if req.LastUpdatedAt != nil {
		updatedAt = timestamppb.New(*req.LastUpdatedAt)
	}

	go func() {
		defer close(ch)
		defer cancel()

		md := metadata.New(map[string]string{"subscriber_id": req.SubscriberID})
		ctx = metadata.NewOutgoingContext(ctx, md)

		stream, errC := s.coreFeed.EventsSubscribe(ctx, &feedproto.EventsSubscribeRequest{
			SubscriberId:      req.SubscriberID,
			SubscriptionTypes: convertTypesToFeedProto(req.SubscriptionTypes),
			LastUpdatedAt:     updatedAt,
		})
		if errC != nil {
			ch <- Result{Err: errC}
			return
		}

		for {
			in, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				ch <- Result{Err: fmt.Errorf("stream.Recv: %w", err)}
			}

			ch <- Result{Item: convertFeedToItem(in)}
		}
	}()

	if slices.Contains(req.SubscriptionTypes, internalproto.FeedItemType_FEED_ITEM_TYPE_VOTE) {
		go func() {
			defer close(ch)
			defer cancel()

			stream, errC := s.coreStorage.VotesSubscribe(ctx, &coreproto.VotesSubscribeRequest{
				LastUpdatedAt: updatedAt,
			})
			if errC != nil {
				ch <- Result{Err: errC}
				return
			}

			for {
				in, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					return
				}

				if err != nil {
					ch <- Result{Err: fmt.Errorf("stream.Recv: %w", err)}
				}

				ch <- Result{Item: convertVoteToInternal(in)}
			}
		}()
	}

	return ch
}

func convertTypesToFeedProto(list []internalproto.FeedItemType) []feedproto.FeedItemType {
	res := make([]feedproto.FeedItemType, 0, len(list))
	for _, item := range list {
		if item == internalproto.FeedItemType_FEED_ITEM_TYPE_VOTE {
			continue
		}

		res = append(res, feedproto.FeedItemType(item))
	}

	return res
}

func convertFeedToItem(in *feedproto.FeedItem) *internalproto.FeedItem {
	if in == nil {
		return nil
	}

	item := &internalproto.FeedItem{
		CreatedAt: in.GetCreatedAt(),
		UpdatedAt: in.GetUpdatedAt(),
		Type:      internalproto.FeedItemType(in.GetType()),
	}

	switch in.GetSnapshot().(type) {
	case *feedproto.FeedItem_Dao:
		item.Snapshot = convertFeedDao(in.GetSnapshot().(*feedproto.FeedItem_Dao))
	case *feedproto.FeedItem_Proposal:
		item.Snapshot = convertFeedProposal(in.GetSnapshot().(*feedproto.FeedItem_Proposal))
	case *feedproto.FeedItem_Delegate:
		item.Snapshot = convertFeedDelegate(in.GetSnapshot().(*feedproto.FeedItem_Delegate))
	}

	return item
}

func convertFeedDao(in *feedproto.FeedItem_Dao) *internalproto.FeedItem_Dao {
	if in == nil || in.Dao == nil {
		return nil
	}

	return &internalproto.FeedItem_Dao{
		Dao: &internalproto.DAO{
			CreatedAt:       in.Dao.GetCreatedAt(),
			InternalId:      in.Dao.GetInternalId(),
			OriginalId:      in.Dao.GetOriginalId(),
			Name:            in.Dao.GetName(),
			Avatar:          in.Dao.GetAvatar(),
			PopularityIndex: in.Dao.GetPopularityIndex(),
			Verified:        in.Dao.GetVerified(),
			Timeline:        convertTimeline(in.Dao.GetTimeline()),
		},
	}
}

func convertFeedProposal(in *feedproto.FeedItem_Proposal) *internalproto.FeedItem_Proposal {
	if in == nil || in.Proposal == nil {
		return nil
	}

	return &internalproto.FeedItem_Proposal{
		Proposal: &internalproto.Proposal{
			CreatedAt:         in.Proposal.GetCreatedAt(),
			Id:                in.Proposal.GetId(),
			DaoInternalId:     in.Proposal.GetDaoInternalId(),
			Author:            in.Proposal.GetAuthor(),
			Title:             in.Proposal.GetTitle(),
			State:             in.Proposal.GetState(),
			Spam:              in.Proposal.GetSpam(),
			Type:              in.Proposal.GetType(),
			Privacy:           in.Proposal.GetPrivacy(),
			Choices:           in.Proposal.GetChoices(),
			OriginalCreatedAt: in.Proposal.GetCreatedAt(),
			VotingStartedAt:   in.Proposal.GetVoteStart(),
			VotingEndedAt:     in.Proposal.GetVoteEnd(),
			Timeline:          convertTimeline(in.Proposal.GetTimeline()),
		},
	}
}

func convertFeedDelegate(in *feedproto.FeedItem_Delegate) *internalproto.FeedItem_Delegate {
	if in == nil || in.Delegate == nil {
		return nil
	}

	return &internalproto.FeedItem_Delegate{
		Delegate: &internalproto.Delegate{
			AddressFrom:   in.Delegate.GetAddressFrom(),
			AddressTo:     in.Delegate.GetAddressTo(),
			DaoInternalId: in.Delegate.GetDaoInternalId(),
			ProposalId:    in.Delegate.GetProposalId(),
			Action:        in.Delegate.GetAction(),
			DueDate:       in.Delegate.GetDueDate(),
		},
	}
}

func convertTimeline(in []*feedproto.Timeline) []*internalproto.Timeline {
	result := make([]*internalproto.Timeline, 0, len(in))
	for _, block := range in {
		result = append(result, &internalproto.Timeline{
			Action:    block.GetAction(),
			CreatedAt: block.GetCreatedAt(),
		})
	}

	return result
}

func convertVoteToInternal(in *coreproto.VoteInfo) *internalproto.FeedItem {
	if in == nil {
		return nil
	}

	return &internalproto.FeedItem{
		CreatedAt: timestamppb.New(time.Unix(int64(in.GetCreated()), 0)),
		Type:      internalproto.FeedItemType_FEED_ITEM_TYPE_VOTE,
		Snapshot: &internalproto.FeedItem_Vote{
			Vote: &internalproto.Vote{
				CreatedAt:     timestamppb.New(time.Unix(int64(in.GetCreated()), 0)),
				DaoInternalId: in.GetDaoId(),
				ProposalId:    in.GetProposalId(),
				VoterAddress:  in.GetVoter(),
				VoteId:        in.GetId(),
				Choice:        in.GetChoice(),
				Reason:        in.GetReason(),
				VotingPower:   in.GetVp(),
			},
		},
	}
}
