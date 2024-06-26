//nolint:dupl
package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionconfig1 "github.com/NpoolPlatform/inspire-gateway/pkg/app/commission/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/commission/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAppCommissionConfig(ctx context.Context, in *npool.UpdateAppCommissionConfigRequest) (*npool.UpdateAppCommissionConfigResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithID(&in.ID, true),
		commissionconfig1.WithEntID(&in.EntID, true),
		commissionconfig1.WithAppID(&in.AppID, true),
		commissionconfig1.WithThresholdAmount(in.ThresholdAmount, false),
		commissionconfig1.WithInvites(in.Invites, false),
		commissionconfig1.WithStartAt(in.StartAt, false),
		commissionconfig1.WithDisabled(in.Disabled, false),
		commissionconfig1.WithLevel(in.Level, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateAppCommissionConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCommission(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateAppCommissionConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppCommissionConfigResponse{
		Info: info,
	}, nil
}

func (s *Server) AdminUpdateAppCommissionConfig(ctx context.Context, in *npool.AdminUpdateAppCommissionConfigRequest) (*npool.AdminUpdateAppCommissionConfigResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithID(&in.ID, true),
		commissionconfig1.WithEntID(&in.EntID, true),
		commissionconfig1.WithAppID(&in.TargetAppID, true),
		commissionconfig1.WithThresholdAmount(in.ThresholdAmount, false),
		commissionconfig1.WithInvites(in.Invites, false),
		commissionconfig1.WithStartAt(in.StartAt, false),
		commissionconfig1.WithDisabled(in.Disabled, false),
		commissionconfig1.WithLevel(in.Level, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateAppCommissionConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCommission(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateAppCommissionConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminUpdateAppCommissionConfigResponse{
		Info: info,
	}, nil
}
