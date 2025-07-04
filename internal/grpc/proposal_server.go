package grpc

import (
	"context"
	"time"

	coredata "github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	internalpb "github.com/goverland-labs/goverland-core-web-api/protocol/storage"
	"go.openly.dev/pointy"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProposalServer struct {
	internalpb.UnimplementedProposalServer

	pc coredata.ProposalClient
}

func NewProposalServer(pc coredata.ProposalClient) *ProposalServer {
	return &ProposalServer{
		pc: pc,
	}
}

func (s *ProposalServer) GetByID(ctx context.Context, req *internalpb.ProposalByIDRequest) (*internalpb.ProposalByIDResponse, error) {
	pr, err := s.pc.GetByID(ctx, &coredata.ProposalByIDRequest{
		ProposalId: req.GetProposalId(),
	})
	if err != nil {
		return nil, err
	}

	return &internalpb.ProposalByIDResponse{
		Proposal: convertProposal(pr.GetProposal()),
	}, nil
}

func convertProposal(pr *coredata.ProposalInfo) *internalpb.ProposalInfo {
	if pr == nil {
		return nil
	}

	timeline := make([]*internalpb.Timeline, 0, len(pr.GetTimeline()))
	for _, info := range pr.GetTimeline() {
		timeline = append(timeline, &internalpb.Timeline{
			Action:    convertAction(info.GetAction()),
			CreatedAt: info.GetCreatedAt(),
		})
	}

	return &internalpb.ProposalInfo{
		Id:                pr.GetId(),
		CreatedAt:         pr.GetCreatedAt(),
		UpdatedAt:         pr.GetUpdatedAt(),
		Author:            pr.GetAuthor(),
		DaoId:             pr.GetDaoId(),
		Title:             pr.GetTitle(),
		State:             pr.GetState(),
		Type:              pr.GetType(),
		Privacy:           pr.GetPrivacy(),
		Timeline:          timeline,
		Spam:              pr.GetSpam(),
		Choices:           pr.GetChoices(),
		OriginalCreatedAt: timestamppb.New(time.Unix(int64(pr.GetCreated()), 0)),
		VotingStartedAt:   timestamppb.New(time.Unix(int64(pr.GetStart()), 0)),
		VotingEndedAt:     timestamppb.New(time.Unix(int64(pr.GetEnd()), 0)),
	}
}

// todo: think how to remove double conversions in core-feed -> core-storage -> core-web-api
func convertAction(action coredata.ProposalTimelineItem_TimelineAction) string {
	switch action {
	case coredata.ProposalTimelineItem_ProposalCreated:
		return "proposal.created"
	case coredata.ProposalTimelineItem_ProposalUpdated:
		return "proposal.updated"
	case coredata.ProposalTimelineItem_ProposalVotingStarted:
		return "proposal.voting.started"
	case coredata.ProposalTimelineItem_ProposalVotingEnded:
		return "proposal.voting.ended"
	case coredata.ProposalTimelineItem_ProposalVotingQuorumReached:
		return "proposal.voting.quorum_reached"
	case coredata.ProposalTimelineItem_ProposalVotingStartsSoon:
		return "proposal.voting.starts_soon"
	case coredata.ProposalTimelineItem_ProposalVotingEndsSoon:
		return "proposal.voting.ends_soon"
	default:
		return ""
	}
}

func (s *ProposalServer) GetByFilter(ctx context.Context, req *internalpb.ProposalByFilterRequest) (*internalpb.ProposalByFilterResponse, error) {
	resp, err := s.pc.GetByFilter(ctx, &coredata.ProposalByFilterRequest{
		Dao:         req.Dao,
		Limit:       req.Limit,
		Offset:      req.Offset,
		ProposalIds: req.GetProposalIds(),
		OnlyActive:  req.OnlyActive,
		Level:       pointy.Pointer(convertReqLevel(req.GetLevel())),
	})
	if err != nil {
		return nil, err
	}

	result := &internalpb.ProposalByFilterResponse{
		Proposals:      make([]*internalpb.ProposalInfo, 0, len(resp.Proposals)),
		TotalCount:     resp.GetTotalCount(),
		ProposalsShort: make([]*internalpb.ProposalShortInfo, 0, len(resp.ProposalsShort)),
	}

	for _, info := range resp.Proposals {
		result.Proposals = append(result.Proposals, convertProposal(info))
	}

	for _, info := range resp.ProposalsShort {
		result.ProposalsShort = append(result.ProposalsShort, convertShortProposal(info))
	}

	return result, nil
}

func convertReqLevel(level internalpb.ProposalInfoLevel) coredata.ProposalInfoLevel {
	switch level {
	case internalpb.ProposalInfoLevel_PROPOSAL_INFO_LEVEL_SHORT:
		return coredata.ProposalInfoLevel_PROPOSAL_INFO_LEVEL_SHORT
	default:
		return coredata.ProposalInfoLevel_PROPOSAL_INFO_LEVEL_FULL
	}
}

func convertShortProposal(pr *coredata.ProposalShortInfo) *internalpb.ProposalShortInfo {
	if pr == nil {
		return nil
	}

	return &internalpb.ProposalShortInfo{
		Id:      pr.GetId(),
		Title:   pr.GetTitle(),
		State:   pr.GetState(),
		Created: pr.GetCreated(),
	}
}
