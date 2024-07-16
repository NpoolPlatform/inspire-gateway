package event

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	event1 "github.com/NpoolPlatform/inspire-gateway/pkg/event"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminUpdateEvent(ctx context.Context, in *npool.AdminUpdateEventRequest) (*npool.AdminUpdateEventResponse, error) {
	handler, err := event1.NewHandler(
		ctx,
		event1.WithID(&in.ID, true),
		event1.WithEntID(&in.EntID, true),
		event1.WithAppID(&in.TargetAppID, true),
		event1.WithCredits(in.Credits, false),
		event1.WithCreditsPerUSD(in.CreditsPerUSD, false),
		event1.WithMaxConsecutive(in.MaxConsecutive, false),
		event1.WithInviterLayers(in.InviterLayers, false),
		event1.WithCouponIDs(in.GetCouponIDs(), false),
		event1.WithCoins(in.GetCoins(), false),
		event1.WithRemoveCouponIDs(in.RemoveCouponIDs, false),
		event1.WithRemoveCoins(in.RemoveCoins, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateEvent",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateEvent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateEvent",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateEventResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminUpdateEventResponse{
		Info: info,
	}, nil
}
