package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionconfig1 "github.com/NpoolPlatform/inspire-gateway/pkg/app/commission/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/commission/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAppCommissionConfigs(ctx context.Context, in *npool.GetAppCommissionConfigsRequest) (*npool.GetAppCommissionConfigsResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.AppID, true),
		commissionconfig1.WithOffset(in.GetOffset()),
		commissionconfig1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCommissionConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCommissions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppCommissionConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCommissionConfigsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppCommissionConfigs(ctx context.Context, in *npool.GetNAppCommissionConfigsRequest) (*npool.GetNAppCommissionConfigsResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.TargetAppID, true),
		commissionconfig1.WithEndAt(in.EndAt, false),
		commissionconfig1.WithOffset(in.GetOffset()),
		commissionconfig1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.GetNAppCommissionConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCommissions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.GetNAppCommissionConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppCommissionConfigsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
