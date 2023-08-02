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
		event1.WithAppID(&in.AppID),
		event1.WithEventType(&in.EventType),
		event1.WithCouponIDs(in.GetCouponIDs()),
		event1.WithCredits(in.Credits),
		event1.WithCreditsPerUSD(in.CreditsPerUSD),
		event1.WithMaxConsecutive(in.MaxConsecutive),
		event1.WithGoodID(in.GoodID),
		event1.WithInviterLayers(in.InviterLayers),
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
