package coin

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coin1 "github.com/NpoolPlatform/inspire-gateway/pkg/event/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminUpdateEventCoin(ctx context.Context, in *npool.AdminUpdateEventCoinRequest) (*npool.AdminUpdateEventCoinResponse, error) {
	handler, err := coin1.NewHandler(
		ctx,
		coin1.WithID(&in.ID, true),
		coin1.WithEntID(&in.EntID, true),
		coin1.WithAppID(&in.TargetAppID, true),
		coin1.WithCoinValue(in.CoinValue, false),
		coin1.WithCoinPreUSD(in.CoinPreUSD, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateEventCoin",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateEventCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateEventCoin(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateEventCoin",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateEventCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminUpdateEventCoinResponse{
		Info: info,
	}, nil
}
