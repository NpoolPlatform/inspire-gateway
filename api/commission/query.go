package commission

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commission1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCommissions(ctx context.Context, in *npool.GetCommissionsRequest) (*npool.GetCommissionsResponse, error) {
	handler, err := commission1.NewHandler(
		ctx,
		commission1.WithAppID(&in.AppID, true),
		commission1.WithUserID(&in.UserID, true),
		commission1.WithOffset(in.GetOffset()),
		commission1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCommission",
			"In", in,
			"Err", err,
		)
		return &npool.GetCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCommissions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCommission",
			"In", in,
			"Err", err,
		)
		return &npool.GetCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCommissionsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppCommissions(ctx context.Context, in *npool.GetAppCommissionsRequest) (*npool.GetAppCommissionsResponse, error) {
	handler, err := commission1.NewHandler(
		ctx,
		commission1.WithAppID(&in.AppID, true),
		commission1.WithEndAt(in.EndAt),
		commission1.WithOffset(in.GetOffset()),
		commission1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCommission",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCommissions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCommission",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCommissionsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
