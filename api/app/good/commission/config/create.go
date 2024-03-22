package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionconfig1 "github.com/NpoolPlatform/inspire-gateway/pkg/app/good/commission/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/good/commission/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateAppGoodCommissionConfig(ctx context.Context, in *npool.CreateAppGoodCommissionConfigRequest) (*npool.CreateAppGoodCommissionConfigResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.AppID, true),
		commissionconfig1.WithAppGoodID(&in.AppGoodID, true),
		commissionconfig1.WithAmountOrPercent(&in.AmountOrPercent, true),
		commissionconfig1.WithStartAt(&in.StartAt, false),
		commissionconfig1.WithInvites(&in.Invites, true),
		commissionconfig1.WithThresholdAmount(&in.ThresholdAmount, true),
		commissionconfig1.WithSettleType(&in.SettleType, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppGoodCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateAppGoodCommissionConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCommissionConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppGoodCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateAppGoodCommissionConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppGoodCommissionConfigResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateNAppGoodCommissionConfig(ctx context.Context, in *npool.CreateNAppGoodCommissionConfigRequest) (*npool.CreateNAppGoodCommissionConfigResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.TargetAppID, true),
		commissionconfig1.WithAppGoodID(&in.AppGoodID, true),
		commissionconfig1.WithAmountOrPercent(&in.AmountOrPercent, true),
		commissionconfig1.WithStartAt(&in.StartAt, false),
		commissionconfig1.WithInvites(&in.Invites, true),
		commissionconfig1.WithThresholdAmount(&in.ThresholdAmount, true),
		commissionconfig1.WithSettleType(&in.SettleType, true),
		commissionconfig1.WithCheckAffiliate(false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateNAppGoodCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateNAppGoodCommissionConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCommissionConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateNAppGoodCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateNAppGoodCommissionConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateNAppGoodCommissionConfigResponse{
		Info: info,
	}, nil
}
