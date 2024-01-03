package cashcontrol

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cashcontrol1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/cashcontrol"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/cashcontrol"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateCashControl(ctx context.Context, in *npool.CreateCashControlRequest) (*npool.CreateCashControlResponse, error) {
	handler, err := cashcontrol1.NewHandler(
		ctx,
		cashcontrol1.WithCouponID(&in.CouponID, true),
		cashcontrol1.WithControlType(&in.ControlType, true),
		cashcontrol1.WithValue(in.Value, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCashControl",
			"In", in,
			"Err", err,
		)
		return &npool.CreateCashControlResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCashControl(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCashControl",
			"In", in,
			"Err", err,
		)
		return &npool.CreateCashControlResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCashControlResponse{
		Info: info,
	}, nil
}
