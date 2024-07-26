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

func (s *Server) CreateEventCoin(ctx context.Context, in *npool.CreateEventCoinRequest) (*npool.CreateEventCoinResponse, error) {
	handler, err := coin1.NewHandler(
		ctx,
		coin1.WithAppID(&in.AppID, true),
		coin1.WithEventID(&in.EventID, true),
		coin1.WithCoinConfigID(&in.CoinConfigID, true),
		coin1.WithCoinValue(&in.CoinValue, true),
		coin1.WithCoinPreUSD(in.CoinPreUSD, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateEventCoin",
			"In", in,
			"Err", err,
		)
		return &npool.CreateEventCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateEvent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateEventCoin",
			"In", in,
			"Err", err,
		)
		return &npool.CreateEventCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateEventCoinResponse{
		Info: info,
	}, nil
}
