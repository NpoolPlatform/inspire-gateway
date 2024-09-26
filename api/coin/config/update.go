package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	config1 "github.com/NpoolPlatform/inspire-gateway/pkg/coin/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminUpdateCoinConfig(ctx context.Context, in *npool.AdminUpdateCoinConfigRequest) (*npool.AdminUpdateCoinConfigResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithID(&in.ID, true),
		config1.WithEntID(&in.EntID, true),
		config1.WithAppID(&in.TargetAppID, true),
		config1.WithMaxValue(in.MaxValue, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateCoinConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateCoinConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCoinConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateCoinConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateCoinConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminUpdateCoinConfigResponse{
		Info: info,
	}, nil
}
