package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	config1 "github.com/NpoolPlatform/inspire-gateway/pkg/task/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminDeleteTaskConfig(ctx context.Context, in *npool.AdminDeleteTaskConfigRequest) (*npool.AdminDeleteTaskConfigResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithID(&in.ID, true),
		config1.WithEntID(&in.EntID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteTaskConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteTaskConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteTaskConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminDeleteTaskConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminDeleteTaskConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminDeleteTaskConfigResponse{
		Info: info,
	}, nil
}
