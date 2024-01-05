package coin

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	couponcoin1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/coin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCouponCoins(ctx context.Context, in *npool.GetCouponCoinsRequest) (*npool.GetCouponCoinsResponse, error) { //nolint
	handler, err := couponcoin1.NewHandler(
		ctx,
		couponcoin1.WithAppID(&in.TargetAppID, true),
		couponcoin1.WithOffset(in.GetOffset()),
		couponcoin1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCouponCoins",
			"In", in,
			"Err", err,
		)
		return &npool.GetCouponCoinsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCouponCoins(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCouponCoins",
			"In", in,
			"Err", err,
		)
		return &npool.GetCouponCoinsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCouponCoinsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppCouponCoins(ctx context.Context, in *npool.GetAppCouponCoinsRequest) (*npool.GetAppCouponCoinsResponse, error) { //nolint
	handler, err := couponcoin1.NewHandler(
		ctx,
		couponcoin1.WithAppID(&in.AppID, true),
		couponcoin1.WithOffset(in.GetOffset()),
		couponcoin1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCouponCoins",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCouponCoinsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCouponCoins(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCouponCoins",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCouponCoinsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCouponCoinsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
