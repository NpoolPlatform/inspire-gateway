package commission

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
		commissionconfig1.WithStartAt(&in.StartAt, false),
		commissionconfig1.WithInvites(&in.Invites, true),
		commissionconfig1.WithThresholdAmount(&in.ThresholdAmount, true),
		commissionconfig1.WithSettleType(&in.SettleType, true),
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

func (s *Server) CreateNAppCommissionConfig(ctx context.Context, in *npool.CreateNAppCommissionConfigRequest) (*npool.CreateNAppCommissionConfigResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.TargetAppID, true),
		commissionconfig1.WithAmountOrPercent(&in.AmountOrPercent, true),
		commissionconfig1.WithStartAt(&in.StartAt, false),
		commissionconfig1.WithInvites(&in.Invites, true),
		commissionconfig1.WithThresholdAmount(&in.ThresholdAmount, true),
		commissionconfig1.WithSettleType(&in.SettleType, true),
		commissionconfig1.WithCheckAffiliate(false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateNAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateNAppCommissionConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCommissionConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateNAppCommissionConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateNAppCommissionConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateNAppCommissionConfigResponse{
		Info: info,
	}, nil
}
