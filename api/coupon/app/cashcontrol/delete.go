package cashcontrol

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cashcontrol1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/cashcontrol"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/cashcontrol"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteCashControl(ctx context.Context, in *npool.DeleteCashControlRequest) (*npool.DeleteCashControlResponse, error) {
	handler, err := cashcontrol1.NewHandler(
		ctx,
		cashcontrol1.WithID(&in.ID, true),
		cashcontrol1.WithEntID(&in.EntID, true),
		cashcontrol1.WithAppID(&in.TargetAppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCashControl",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCashControlResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteCashControl(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCashControl",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCashControlResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteCashControlResponse{
		Info: info,
	}, nil
}
