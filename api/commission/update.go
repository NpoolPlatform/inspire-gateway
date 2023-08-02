package commission

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commission1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateCommission(ctx context.Context, in *npool.UpdateCommissionRequest) (*npool.UpdateCommissionResponse, error) {
	handler, err := commission1.NewHandler(
		ctx,
		commission1.WithID(&in.ID),
		commission1.WithAppID(&in.AppID),
		commission1.WithAmountOrPercent(in.AmountOrPercent),
		commission1.WithStartAt(in.StartAt),
		commission1.WithThreshold(in.Threshold),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCommission",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCommission(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCommission",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCommissionResponse{
		Info: info,
	}, nil
}
