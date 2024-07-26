package coupon

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coupon1 "github.com/NpoolPlatform/inspire-gateway/pkg/event/coupon"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coupon"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminDeleteEventCoupon(ctx context.Context, in *npool.AdminDeleteEventCouponRequest) (*npool.AdminDeleteEventCouponResponse, error) {
	handler, err := coupon1.NewHandler(
		ctx,
		coupon1.WithID(&in.ID, true),
		coupon1.WithEntID(&in.EntID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteEventCoupon",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteEventCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteEventCoupon(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteEventCoupon",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteEventCouponResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminDeleteEventCouponResponse{
		Info: info,
	}, nil
}
