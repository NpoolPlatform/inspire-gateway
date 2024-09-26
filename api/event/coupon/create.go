//nolint:dupl
package coupon

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coupon1 "github.com/NpoolPlatform/inspire-gateway/pkg/event/coupon"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coupon"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateEventCoupon(ctx context.Context, in *npool.CreateEventCouponRequest) (*npool.CreateEventCouponResponse, error) {
	handler, err := coupon1.NewHandler(
		ctx,
		coupon1.WithAppID(&in.AppID, true),
		coupon1.WithEventID(&in.EventID, true),
		coupon1.WithCouponID(&in.CouponID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateEventCoupon",
			"In", in,
			"Err", err,
		)
		return &npool.CreateEventCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateEvent(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateEventCoupon",
			"In", in,
			"Err", err,
		)
		return &npool.CreateEventCouponResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateEventCouponResponse{
		Info: info,
	}, nil
}
