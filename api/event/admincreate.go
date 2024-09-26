//nolint:dupl
package event

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	event1 "github.com/NpoolPlatform/inspire-gateway/pkg/event"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminCreateEvent(ctx context.Context, in *npool.AdminCreateEventRequest) (*npool.AdminCreateEventResponse, error) {
	handler, err := event1.NewHandler(
		ctx,
		event1.WithAppID(&in.TargetAppID, true),
		event1.WithEventType(&in.EventType, true),
		event1.WithCredits(in.Credits, false),
		event1.WithCreditsPerUSD(in.CreditsPerUSD, false),
		event1.WithAppGoodID(in.AppGoodID, false),
		event1.WithMaxConsecutive(in.MaxConsecutive, false),
		event1.WithInviterLayers(in.InviterLayers, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateEvent",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateEvent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateEvent",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateEventResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminCreateEventResponse{
		Info: info,
	}, nil
}
