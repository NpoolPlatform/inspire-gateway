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

func (s *Server) AdminDeleteEvent(ctx context.Context, in *npool.AdminDeleteEventRequest) (*npool.AdminDeleteEventResponse, error) {
	handler, err := event1.NewHandler(
		ctx,
		event1.WithID(&in.ID, true),
		event1.WithEntID(&in.EntID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteEvent",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateEvent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteEvent",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteEventResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminDeleteEventResponse{
		Info: info,
	}, nil
}
