//nolint:dupl
package coin

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coin1 "github.com/NpoolPlatform/inspire-gateway/pkg/event/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminCreateEventCoin(ctx context.Context, in *npool.AdminCreateEventCoinRequest) (*npool.AdminCreateEventCoinResponse, error) {
	handler, err := coin1.NewHandler(
		ctx,
		coin1.WithAppID(&in.TargetAppID, true),
		coin1.WithEventID(&in.EventID, true),
		coin1.WithCoinConfigID(&in.CoinConfigID, true),
		coin1.WithCoinValue(&in.CoinValue, true),
		coin1.WithCoinPerUSD(in.CoinPerUSD, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateEventCoin",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateEventCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateEvent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateEventCoin",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateEventCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminCreateEventCoinResponse{
		Info: info,
	}, nil
}
