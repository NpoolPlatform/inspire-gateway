package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	config1 "github.com/NpoolPlatform/inspire-gateway/pkg/coin/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminCreateCoinConfig(ctx context.Context, in *npool.AdminCreateCoinConfigRequest) (*npool.AdminCreateCoinConfigResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithAppID(&in.TargetAppID, true),
		config1.WithCoinTypeID(&in.CoinTypeID, true),
		config1.WithMaxValue(&in.MaxValue, true),
		config1.WithAllocated(&in.Allocated, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateCoinConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateCoinConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCoinConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateCoinConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateCoinConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminCreateCoinConfigResponse{
		Info: info,
	}, nil
}
