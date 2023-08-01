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
		commission1.WithAppID(&in.AppID),
		commission1.WithUserID(&in.UserID),
		commission1.WithTargetUserID(&in.TargetUserID),
		commission1.WithGoodID(&in.GoodID),
		commission1.WithSettleType(&in.SettleType),
		commission1.WithAmountOrPercent(&in.AmountOrPercent),
		commission1.WithStartAt(&in.StartAt),
		commission1.WithSettleMode(&in.SettleMode),
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
		commission1.WithAppID(&in.AppID),
		commission1.WithTargetUserID(&in.TargetUserID),
		commission1.WithGoodID(&in.GoodID),
		commission1.WithSettleType(&in.SettleType),
		commission1.WithAmountOrPercent(&in.AmountOrPercent),
		commission1.WithStartAt(&in.StartAt),
		commission1.WithSettleMode(&in.SettleMode),
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
