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
		coupon1.WithID(&in.ID),
		coupon1.WithAppID(&in.TargetAppID),
		coupon1.WithDenomination(in.Denomination),
		coupon1.WithCirculation(in.Circulation),
		coupon1.WithStartAt(in.StartAt),
		coupon1.WithDurationDays(in.DurationDays),
		coupon1.WithMessage(in.Message),
		coupon1.WithName(in.Name),
		coupon1.WithUserID(in.TargetUserID),
		coupon1.WithGoodID(in.GoodID),
		coupon1.WithThreshold(in.Threshold),
		coupon1.WithCouponConstraint(in.CouponConstraint),
		coupon1.WithRandom(in.Random),
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
