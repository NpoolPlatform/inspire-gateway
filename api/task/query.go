package task

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	task1 "github.com/NpoolPlatform/inspire-gateway/pkg/task"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//nolint:dupl
func (s *Server) AdminGetTasks(ctx context.Context, in *npool.AdminGetTasksRequest) (*npool.AdminGetTasksResponse, error) {
	handler, err := task1.NewHandler(
		ctx,
		task1.WithAppID(&in.TargetAppID, true),
		task1.WithUserID(&in.TargetUserID, true),
		task1.WithOffset(in.GetOffset()),
		task1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetTasks",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetTasksResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetUserTasks(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetTasks",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetTasksResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetTasksResponse{
		Infos: infos,
		Total: total,
	}, nil
}

//nolint:dupl
func (s *Server) GetMyTasks(ctx context.Context, in *npool.GetMyTasksRequest) (*npool.GetMyTasksResponse, error) {
	handler, err := task1.NewHandler(
		ctx,
		task1.WithAppID(&in.AppID, true),
		task1.WithUserID(&in.UserID, true),
		task1.WithOffset(in.GetOffset()),
		task1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMyTasks",
			"In", in,
			"Err", err,
		)
		return &npool.GetMyTasksResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetUserTasks(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMyTasks",
			"In", in,
			"Err", err,
		)
		return &npool.GetMyTasksResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetMyTasksResponse{
		Infos: infos,
		Total: total,
	}, nil
}
