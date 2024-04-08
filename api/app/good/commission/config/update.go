//nolint:dupl
package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionconfig1 "github.com/NpoolPlatform/inspire-gateway/pkg/app/good/commission/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/good/commission/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAppGoodCommissionConfig(ctx context.Context, in *npool.UpdateAppGoodCommissionConfigRequest) (*npool.UpdateAppGoodCommissionConfigResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithID(&in.ID, true),
		commissionconfig1.WithEntID(&in.EntID, true),
		commissionconfig1.WithAppID(&in.AppID, true),
		commissionconfig1.WithThresholdAmount(in.ThresholdAmount, false),
		commissionconfig1.WithInvites(in.Invites, false),
		commissionconfig1.WithStartAt(in.StartAt, false),
		commissionconfig1.WithDisabled(in.Disabled, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppGoodCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateAppGoodCommissionConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCommission(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppGoodCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateAppGoodCommissionConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppGoodCommissionConfigResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateNAppGoodCommissionConfig(ctx context.Context, in *npool.UpdateNAppGoodCommissionConfigRequest) (*npool.UpdateNAppGoodCommissionConfigResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithID(&in.ID, true),
		commissionconfig1.WithEntID(&in.EntID, true),
		commissionconfig1.WithAppID(&in.TargetAppID, true),
		commissionconfig1.WithThresholdAmount(in.ThresholdAmount, false),
		commissionconfig1.WithInvites(in.Invites, false),
		commissionconfig1.WithStartAt(in.StartAt, false),
		commissionconfig1.WithDisabled(in.Disabled, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateNAppGoodCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateNAppGoodCommissionConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCommission(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateNAppGoodCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateNAppGoodCommissionConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateNAppGoodCommissionConfigResponse{
		Info: info,
	}, nil
}
