package cashcontrol

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cashcontrol1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/cashcontrol"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/cashcontrol"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateCashControl(ctx context.Context, in *npool.UpdateCashControlRequest) (*npool.UpdateCashControlResponse, error) {
	handler, err := cashcontrol1.NewHandler(
		ctx,
		cashcontrol1.WithID(&in.ID, true),
		cashcontrol1.WithEntID(&in.EntID, true),
		cashcontrol1.WithAppID(&in.TargetAppID, true),
		cashcontrol1.WithValue(in.Value, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCashControl",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateCashControlResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCashControl(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCashControl",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateCashControlResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCashControlResponse{
		Info: info,
	}, nil
}
