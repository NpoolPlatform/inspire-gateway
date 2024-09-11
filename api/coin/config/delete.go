package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	config1 "github.com/NpoolPlatform/inspire-gateway/pkg/coin/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminDeleteCoinConfig(ctx context.Context, in *npool.AdminDeleteCoinConfigRequest) (*npool.AdminDeleteCoinConfigResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithID(&in.ID, true),
		config1.WithEntID(&in.EntID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteCoinConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteCoinConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteCoinConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteCoinConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteCoinConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminDeleteCoinConfigResponse{
		Info: info,
	}, nil
}
