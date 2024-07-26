package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	config1 "github.com/NpoolPlatform/inspire-gateway/pkg/task/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminGetTaskConfigs(ctx context.Context, in *npool.AdminGetTaskConfigsRequest) (*npool.AdminGetTaskConfigsResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithAppID(&in.TargetAppID, true),
		config1.WithOffset(in.GetOffset()),
		config1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetTaskConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetTaskConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetTaskConfigs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetTaskConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetTaskConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetTaskConfigsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
