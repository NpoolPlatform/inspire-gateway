package commission

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commission1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateCommission(ctx context.Context, in *npool.CreateCommissionRequest) (*npool.CreateCommissionResponse, error) {
	handler, err := commission1.NewHandler(
		ctx,
		commission1.WithAppID(&in.AppID, true),
		commission1.WithUserID(&in.UserID, true),
		commission1.WithTargetUserID(&in.TargetUserID, true),
		commission1.WithAppGoodID(&in.AppGoodID, true),
		commission1.WithSettleType(&in.SettleType, true),
		commission1.WithSettleAmountType(&in.SettleAmountType, true),
		commission1.WithAmountOrPercent(&in.AmountOrPercent, true),
		commission1.WithSettleMode(&in.SettleMode, true),
		commission1.WithStartAt(&in.StartAt, false),
		commission1.WithSettleInterval(&in.SettleInterval, true),
		commission1.WithThreshold(in.Threshold, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCommission",
			"In", in,
			"Err", err,
		)
		return &npool.CreateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCommission(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCommission",
			"In", in,
			"Err", err,
		)
		return &npool.CreateCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCommissionResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateUserCommission(ctx context.Context, in *npool.CreateUserCommissionRequest) (*npool.CreateUserCommissionResponse, error) {
	handler, err := commission1.NewHandler(
		ctx,
		commission1.WithAppID(&in.AppID, true),
		commission1.WithTargetUserID(&in.TargetUserID, true),
		commission1.WithAppGoodID(&in.AppGoodID, true),
		commission1.WithSettleType(&in.SettleType, true),
		commission1.WithSettleAmountType(&in.SettleAmountType, true),
		commission1.WithAmountOrPercent(&in.AmountOrPercent, true),
		commission1.WithStartAt(&in.StartAt, true),
		commission1.WithSettleMode(&in.SettleMode, true),
		commission1.WithSettleInterval(&in.SettleInterval, true),
		commission1.WithThreshold(in.Threshold, false),
		commission1.WithCheckAffiliate(false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateUserCommission",
			"In", in,
			"Err", err,
		)
		return &npool.CreateUserCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCommission(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateUserCommission",
			"In", in,
			"Err", err,
		)
		return &npool.CreateUserCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateUserCommissionResponse{
		Info: info,
	}, nil
}
