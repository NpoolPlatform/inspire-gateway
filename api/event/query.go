package event

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"

	event1 "github.com/NpoolPlatform/inspire-gateway/pkg/event"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetEvents(ctx context.Context, in *npool.GetEventsRequest) (*npool.GetEventsResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetEvents", "AppID", in.GetAppID(), "Error", err)
		return &npool.GetEventsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := constant.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	conds := &mgrpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: in.GetAppID()},
	}

	infos, total, err := event1.GetEvents(ctx, conds, in.GetOffset(), limit)
	if err != nil {
		logger.Sugar().Errorw("GetEvents", "Conds", conds, "Error", err)
		return &npool.GetEventsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetEventsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
