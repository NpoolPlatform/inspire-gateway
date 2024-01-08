package coin

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	couponcoin1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/coin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateCouponCoin(ctx context.Context, in *npool.CreateCouponCoinRequest) (*npool.CreateCouponCoinResponse, error) {
	handler, err := couponcoin1.NewHandler(
		ctx,
		couponcoin1.WithAppID(&in.TargetAppID, true),
		couponcoin1.WithCoinTypeID(&in.CoinTypeID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCouponCoin",
			"In", in,
			"Err", err,
		)
		return &npool.CreateCouponCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCouponCoin(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCouponCoin",
			"In", in,
			"Err", err,
		)
		return &npool.CreateCouponCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCouponCoinResponse{
		Info: info,
	}, nil
}
