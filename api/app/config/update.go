//nolint:dupl
package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	appconfig1 "github.com/NpoolPlatform/inspire-gateway/pkg/app/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAppConfig(ctx context.Context, in *npool.UpdateAppConfigRequest) (*npool.UpdateAppConfigResponse, error) {
	handler, err := appconfig1.NewHandler(
		ctx,
		appconfig1.WithID(&in.ID, true),
		appconfig1.WithEntID(&in.EntID, true),
		appconfig1.WithAppID(&in.AppID, true),
		appconfig1.WithStartAt(in.StartAt, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateAppConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateAppConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateAppConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppConfigResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateNAppConfig(ctx context.Context, in *npool.UpdateNAppConfigRequest) (*npool.UpdateNAppConfigResponse, error) {
	handler, err := appconfig1.NewHandler(
		ctx,
		appconfig1.WithID(&in.ID, true),
		appconfig1.WithEntID(&in.EntID, true),
		appconfig1.WithAppID(&in.TargetAppID, true),
		appconfig1.WithStartAt(in.StartAt, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateNAppConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateNAppConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateAppConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateNAppConfig",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateNAppConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateNAppConfigResponse{
		Info: info,
	}, nil
}
