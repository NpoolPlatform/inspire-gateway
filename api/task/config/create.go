package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	config1 "github.com/NpoolPlatform/inspire-gateway/pkg/task/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminCreateTaskConfig(ctx context.Context, in *npool.AdminCreateTaskConfigRequest) (*npool.AdminCreateTaskConfigResponse, error) {
	handler, err := config1.NewHandler(
		ctx,
		config1.WithAppID(&in.TargetAppID, true),
		config1.WithEventID(&in.EventID, true),
		config1.WithTaskType(&in.TaskType, true),
		config1.WithName(&in.Name, true),
		config1.WithTaskDesc(&in.TaskDesc, true),
		config1.WithStepGuide(&in.StepGuide, true),
		config1.WithRecommendMessage(&in.RecommendMessage, true),
		config1.WithIndex(&in.Index, true),
		config1.WithMaxRewardCount(&in.MaxRewardCount, true),
		config1.WithCooldownSecord(&in.CooldownSecord, true),
		config1.WithLastTaskID(in.LastTaskID, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateTaskConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateTaskConfigResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateTaskConfig(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminCreateTaskConfig",
			"In", in,
			"Err", err,
		)
		return &npool.AdminCreateTaskConfigResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminCreateTaskConfigResponse{
		Info: info,
	}, nil
}
