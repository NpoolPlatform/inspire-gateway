package coupon

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coupon1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//nolint
func (s *Server) CreateCoupon(ctx context.Context, in *npool.CreateCouponRequest) (*npool.CreateCouponResponse, error) {
	handler, err := coupon1.NewHandler(
		ctx,
		coupon1.WithIssuedBy(&in.UserID),
		coupon1.WithAppID(&in.TargetAppID),
		coupon1.WithCouponType(&in.CouponType),
		coupon1.WithDenomination(&in.Denomination),
		coupon1.WithCirculation(&in.Circulation),
		coupon1.WithStartAt(&in.StartAt),
		coupon1.WithDurationDays(&in.DurationDays),
		coupon1.WithMessage(&in.Message),
		coupon1.WithName(&in.Name),
		coupon1.WithUserID(in.TargetUserID),
		coupon1.WithGoodID(in.GoodID),
		coupon1.WithThreshold(in.Threshold),
		coupon1.WithCouponConstraint(&in.CouponConstraint),
		coupon1.WithRandom(&in.Random),
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
