package commission

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	appconfig1 "github.com/NpoolPlatform/inspire-gateway/pkg/app/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateAppConfig(ctx context.Context, in *npool.CreateAppConfigRequest) (*npool.CreateAppConfigResponse, error) {
	handler, err := appconfig1.NewHandler(
		ctx,
		appconfig1.WithAppID(&in.AppID, true),
		appconfig1.WithSettleMode(&in.SettleMode, true),
		appconfig1.WithSettleAmountType(&in.SettleAmountType, true),
		appconfig1.WithSettleInterval(&in.SettleInterval, true),
		appconfig1.WithCommissionType(&in.CommissionType, true),
		appconfig1.WithSettleBenefit(&in.SettleBenefit, true),
		appconfig1.WithStartAt(&in.StartAt, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateAppConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateAppConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateAppConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppConfigResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateNAppConfig(ctx context.Context, in *npool.CreateNAppConfigRequest) (*npool.CreateNAppConfigResponse, error) {
	handler, err := appconfig1.NewHandler(
		ctx,
		appconfig1.WithAppID(&in.TargetAppID, true),
		appconfig1.WithSettleMode(&in.SettleMode, true),
		appconfig1.WithSettleAmountType(&in.SettleAmountType, true),
		appconfig1.WithSettleInterval(&in.SettleInterval, true),
		appconfig1.WithCommissionType(&in.CommissionType, true),
		appconfig1.WithSettleBenefit(&in.SettleBenefit, true),
		appconfig1.WithStartAt(&in.StartAt, false),
		appconfig1.WithCheckAffiliate(false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateNAppConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateNAppConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateAppConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateNAppConfig",
			"In", in,
			"Err", err,
		)
		return &npool.CreateNAppConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateNAppConfigResponse{
		Info: info,
	}, nil
}
