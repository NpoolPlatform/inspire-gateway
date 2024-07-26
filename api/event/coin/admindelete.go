package coin

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coin1 "github.com/NpoolPlatform/inspire-gateway/pkg/event/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminDeleteEventCoin(ctx context.Context, in *npool.AdminDeleteEventCoinRequest) (*npool.AdminDeleteEventCoinResponse, error) {
	handler, err := coin1.NewHandler(
		ctx,
		coin1.WithID(&in.ID, true),
		coin1.WithEntID(&in.EntID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteEventCoin",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteEventCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateEventCoin(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteEventCoin",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteEventCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminDeleteEventCoinResponse{
		Info: info,
	}, nil
}
