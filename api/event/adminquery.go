package event

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	event1 "github.com/NpoolPlatform/inspire-gateway/pkg/event"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminGetEvents(ctx context.Context, in *npool.AdminGetEventsRequest) (*npool.AdminGetEventsResponse, error) {
	handler, err := event1.NewHandler(
		ctx,
		event1.WithAppID(&in.TargetAppID, true),
		event1.WithOffset(in.GetOffset()),
		event1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetEvents",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetEventsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetEvents(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetEvents",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetEventsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetEventsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
