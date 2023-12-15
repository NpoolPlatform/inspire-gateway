package coupon

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coupon1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateCoupon(ctx context.Context, in *npool.CreateCouponRequest) (*npool.CreateCouponResponse, error) {
	handler, err := coupon1.NewHandler(
		ctx,
		coupon1.WithAppID(&in.TargetAppID, true),
		coupon1.WithName(&in.Name, true),
		coupon1.WithMessage(&in.Message, true),
		coupon1.WithCouponType(&in.CouponType, true),
		coupon1.WithDenomination(&in.Denomination, true),
		coupon1.WithCirculation(&in.Circulation, true),
		coupon1.WithStartAt(&in.StartAt, true),
		coupon1.WithDurationDays(&in.DurationDays, true),
		coupon1.WithCouponConstraint(&in.CouponConstraint, true),
		coupon1.WithCouponScope(&in.CouponScope, true),
		coupon1.WithIssuedBy(&in.UserID, true),
		coupon1.WithRandom(&in.Random, false),
		coupon1.WithUserID(in.TargetUserID, false),
		coupon1.WithThreshold(in.Threshold, false),
		coupon1.WithCashableProbabilityPerMillion(in.CashableProbabilityPerMillion, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCoupon",
			"In", in,
			"Err", err,
		)
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCoupon(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCoupon",
			"In", in,
			"Err", err,
		)
		return &npool.CreateCouponResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCouponResponse{
		Info: info,
	}, nil
}
