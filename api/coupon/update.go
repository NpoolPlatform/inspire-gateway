package coupon

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coupon1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateCoupon(ctx context.Context, in *npool.UpdateCouponRequest) (*npool.UpdateCouponResponse, error) {
	handler, err := coupon1.NewHandler(
		ctx,
		coupon1.WithID(&in.ID, true),
		coupon1.WithEntID(&in.EntID, true),
		coupon1.WithAppID(&in.TargetAppID, true),
		coupon1.WithDenomination(in.Denomination, false),
		coupon1.WithCirculation(in.Circulation, false),
		coupon1.WithStartAt(in.StartAt, false),
		coupon1.WithDurationDays(in.DurationDays, false),
		coupon1.WithMessage(in.Message, false),
		coupon1.WithName(in.Name, false),
		coupon1.WithUserID(in.TargetUserID, false),
		coupon1.WithThreshold(in.Threshold, false),
		coupon1.WithCouponConstraint(in.CouponConstraint, false),
		coupon1.WithRandom(in.Random, false),
		coupon1.WithCouponScope(in.CouponScope, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCoupon",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCoupon(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCoupon",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateCouponResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCouponResponse{
		Info: info,
	}, nil
}
