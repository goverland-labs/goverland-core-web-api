package grpc

import (
	"context"

	coredata "github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	internalpb "github.com/goverland-labs/goverland-core-web-api/protocol/storage"
)

type DaoServer struct {
	internalpb.UnimplementedDaoServer

	dc coredata.DaoClient
}

func NewDaoServer(dc coredata.DaoClient) *DaoServer {
	return &DaoServer{
		dc: dc,
	}
}

func (s *DaoServer) GetByID(ctx context.Context, req *internalpb.DaoByIDRequest) (*internalpb.DaoByIDResponse, error) {
	data, err := s.dc.GetByID(ctx, &coredata.DaoByIDRequest{
		DaoId: req.GetDaoId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &internalpb.DaoByIDResponse{
		Dao: &internalpb.DaoInfo{
			Id:              data.GetDao().GetId(),
			CreatedAt:       data.GetDao().GetCreatedAt(),
			UpdatedAt:       data.GetDao().GetUpdatedAt(),
			Name:            data.GetDao().GetName(),
			Avatar:          data.GetDao().GetAvatar(),
			Alias:           data.GetDao().GetAlias(),
			Verified:        data.GetDao().GetVerified(),
			PopularityIndex: data.GetDao().GetPopularityIndex(),
		},
	}, nil
}
