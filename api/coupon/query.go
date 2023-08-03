//nolint:dupl
package coupon

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coupon1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCoupons(ctx context.Context, in *npool.GetCouponsRequest) (*npool.GetCouponsResponse, error) {
	handler, err := coupon1.NewHandler(
		ctx,
		coupon1.WithAppID(&in.AppID),
		coupon1.WithCouponType(in.CouponType),
		coupon1.WithOffset(in.GetOffset()),
		coupon1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCoupons",
			"In", in,
			"Err", err,
		)
		return &npool.GetCouponsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoupons(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCoupons",
			"In", in,
			"Err", err,
		)
		return &npool.GetCouponsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCouponsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppCoupons(ctx context.Context, in *npool.GetAppCouponsRequest) (*npool.GetAppCouponsResponse, error) {
	handler, err := coupon1.NewHandler(
		ctx,
		coupon1.WithAppID(&in.TargetAppID),
		coupon1.WithOffset(in.GetOffset()),
		coupon1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCoupons",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCouponsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoupons(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCoupons",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCouponsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCouponsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
