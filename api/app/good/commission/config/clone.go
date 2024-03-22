//nolint:dupl
package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	config1 "github.com/NpoolPlatform/inspire-gateway/pkg/app/good/commission/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/good/commission/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CloneAppGoodCommissionConfigs(ctx context.Context, in *npool.CloneAppGoodCommissionConfigsRequest) (*npool.CloneAppGoodCommissionConfigsResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithAppID(&in.AppID, true),
		config1.WithFromAppGoodID(&in.FromAppGoodID, true),
		config1.WithToAppGoodID(&in.ToAppGoodID, true),
		config1.WithScalePercent(&in.ScalePercent, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CloneAppGoodCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.CloneAppGoodCommissionConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = handler.CloneCommissionConfigs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CloneAppGoodCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.CloneAppGoodCommissionConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CloneAppGoodCommissionConfigsResponse{}, nil
}

func (s *Server) CloneNAppGoodCommissionConfigs(ctx context.Context, in *npool.CloneNAppGoodCommissionConfigsRequest) (*npool.CloneNAppGoodCommissionConfigsResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithAppID(&in.TargetAppID, true),
		config1.WithFromAppGoodID(&in.FromAppGoodID, true),
		config1.WithToAppGoodID(&in.ToAppGoodID, true),
		config1.WithScalePercent(&in.ScalePercent, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CloneNAppGoodCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.CloneNAppGoodCommissionConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = handler.CloneCommissionConfigs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CloneNAppGoodCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.CloneNAppGoodCommissionConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CloneNAppGoodCommissionConfigsResponse{}, nil
}
