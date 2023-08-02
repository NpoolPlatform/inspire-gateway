package allocated

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/allocated"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/allocated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCoupons(ctx context.Context, in *npool.GetCouponsRequest) (*npool.GetCouponsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.AppID),
		allocated1.WithUserID(&in.UserID),
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
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.AppID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppCoupons",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCouponsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoupons(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppCoupons",
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
