package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	config1 "github.com/NpoolPlatform/inspire-gateway/pkg/task/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminUpdateTaskConfig(ctx context.Context, in *npool.AdminUpdateTaskConfigRequest) (*npool.AdminUpdateTaskConfigResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithID(&in.ID, true),
		config1.WithEntID(&in.EntID, true),
		config1.WithAppID(&in.TargetAppID, true),
		config1.WithTaskType(in.TaskType, false),
		config1.WithName(in.Name, false),
		config1.WithTaskDesc(in.TaskDesc, false),
		config1.WithStepGuide(in.StepGuide, false),
		config1.WithRecommendMessage(in.RecommendMessage, false),
		config1.WithIndex(in.Index, false),
		config1.WithMaxRewardCount(in.MaxRewardCount, false),
		config1.WithCooldownSecord(in.CooldownSecord, false),
		config1.WithLastTaskID(in.LastTaskID, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateTaskConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateTaskConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateTaskConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminUpdateTaskConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminUpdateTaskConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminUpdateTaskConfigResponse{
		Info: info,
	}, nil
}
