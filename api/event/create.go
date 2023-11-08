package event

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	event1 "github.com/NpoolPlatform/inspire-gateway/pkg/event"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateEvent(ctx context.Context, in *npool.CreateEventRequest) (*npool.CreateEventResponse, error) {
	handler, err := event1.NewHandler(
		ctx,
		event1.WithAppID(&in.AppID, true),
		event1.WithEventType(&in.EventType, true),
		event1.WithCredits(in.Credits, false),
		event1.WithCreditsPerUSD(in.CreditsPerUSD, false),
		event1.WithAppGoodID(in.AppGoodID, false),
		event1.WithMaxConsecutive(in.MaxConsecutive, false),
		event1.WithInviterLayers(in.InviterLayers, false),
		event1.WithCouponIDs(in.GetCouponIDs()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateEvent",
			"In", in,
			"Err", err,
		)
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateEvent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateEvent",
			"In", in,
			"Err", err,
		)
		return &npool.CreateEventResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateEventResponse{
		Info: info,
	}, nil
}
