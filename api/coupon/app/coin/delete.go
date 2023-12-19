//nolint
package coin

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	couponcoin1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/coin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteCouponCoin(ctx context.Context, in *npool.DeleteCouponCoinRequest) (*npool.DeleteCouponCoinResponse, error) {
	handler, err := couponcoin1.NewHandler(
		ctx,
		couponcoin1.WithID(&in.ID, true),
		couponcoin1.WithEntID(&in.EntID, true),
		couponcoin1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCouponCoin",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCouponCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteCouponCoin(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCouponCoin",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCouponCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteCouponCoinResponse{
		Info: info,
	}, nil
}
