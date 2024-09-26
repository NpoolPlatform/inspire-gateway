package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	config1 "github.com/NpoolPlatform/inspire-gateway/pkg/coin/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminGetCoinConfigs(ctx context.Context, in *npool.AdminGetCoinConfigsRequest) (*npool.AdminGetCoinConfigsResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithAppID(&in.TargetAppID, true),
		config1.WithOffset(in.GetOffset()),
		config1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetCoinConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetCoinConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoinConfigs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetCoinConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetCoinConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetCoinConfigsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
