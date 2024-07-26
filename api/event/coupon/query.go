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

func (s *Server) GetEventCoupons(ctx context.Context, in *npool.GetEventCouponsRequest) (*npool.GetEventCouponsResponse, error) {
	handler, err := coupon1.NewHandler(
		ctx,
		coupon1.WithAppID(&in.AppID, true),
		coupon1.WithOffset(in.GetOffset()),
		coupon1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetEventCoupons",
			"In", in,
			"Err", err,
		)
		return &npool.GetEventCouponsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetEventCoupons(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetEventCoupons",
			"In", in,
			"Err", err,
		)
		return &npool.GetEventCouponsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetEventCouponsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
