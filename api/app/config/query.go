//nolint:dupl
package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionconfig1 "github.com/NpoolPlatform/inspire-gateway/pkg/app/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAppConfigs(ctx context.Context, in *npool.GetAppConfigsRequest) (*npool.GetAppConfigsResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.AppID, true),
		commissionconfig1.WithEndAt(in.EndAt, false),
		commissionconfig1.WithOffset(in.GetOffset()),
		commissionconfig1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAppConfigs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppConfigsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) AdminGetAppConfigs(ctx context.Context, in *npool.AdminGetAppConfigsRequest) (*npool.AdminGetAppConfigsResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.TargetAppID, true),
		commissionconfig1.WithEndAt(in.EndAt, false),
		commissionconfig1.WithOffset(in.GetOffset()),
		commissionconfig1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetAppConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetAppConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAppConfigs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetAppConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetAppConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetAppConfigsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
