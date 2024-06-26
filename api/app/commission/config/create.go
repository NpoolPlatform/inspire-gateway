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

func (s *Server) CreateAppCommissionConfig(ctx context.Context, in *npool.CreateAppCommissionConfigRequest) (*npool.CreateAppCommissionConfigResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.AppID, true),
		commissionconfig1.WithAmountOrPercent(&in.AmountOrPercent, true),
		commissionconfig1.WithStartAt(in.StartAt, false),
		commissionconfig1.WithInvites(&in.Invites, true),
		commissionconfig1.WithThresholdAmount(&in.ThresholdAmount, true),
		commissionconfig1.WithSettleType(&in.SettleType, true),
		commissionconfig1.WithDisabled(&in.Disabled, false),
		commissionconfig1.WithLevel(&in.Level, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateAppCommissionConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCommissionConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateAppCommissionConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppCommissionConfigResponse{
		Info: info,
	}, nil
}

func (s *Server) AdminCreateAppCommissionConfig(ctx context.Context, in *npool.AdminCreateAppCommissionConfigRequest) (*npool.AdminCreateAppCommissionConfigResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.TargetAppID, true),
		commissionconfig1.WithAmountOrPercent(&in.AmountOrPercent, true),
		commissionconfig1.WithStartAt(in.StartAt, false),
		commissionconfig1.WithInvites(&in.Invites, true),
		commissionconfig1.WithThresholdAmount(&in.ThresholdAmount, true),
		commissionconfig1.WithSettleType(&in.SettleType, true),
		commissionconfig1.WithDisabled(&in.Disabled, false),
		commissionconfig1.WithLevel(&in.Level, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateAppCommissionConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCommissionConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateAppCommissionConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminCreateAppCommissionConfigResponse{
		Info: info,
	}, nil
}
