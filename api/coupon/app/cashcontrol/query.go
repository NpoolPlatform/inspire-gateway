package cashcontrol

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cashcontrol1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/cashcontrol"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/cashcontrol"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCashControls(ctx context.Context, in *npool.GetCashControlsRequest) (*npool.GetCashControlsResponse, error) { //nolint
	handler, err := cashcontrol1.NewHandler(
		ctx,
		cashcontrol1.WithAppID(&in.TargetAppID, true),
		cashcontrol1.WithOffset(in.GetOffset()),
		cashcontrol1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCashControls",
			"In", in,
			"Err", err,
		)
		return &npool.GetCashControlsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCashControls(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCashControls",
			"In", in,
			"Err", err,
		)
		return &npool.GetCashControlsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCashControlsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppCashControls(ctx context.Context, in *npool.GetAppCashControlsRequest) (*npool.GetAppCashControlsResponse, error) { //nolint
	handler, err := cashcontrol1.NewHandler(
		ctx,
		cashcontrol1.WithAppID(&in.AppID, true),
		cashcontrol1.WithOffset(in.GetOffset()),
		cashcontrol1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCashControls",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCashControlsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCashControls(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCashControls",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCashControlsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCashControlsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
