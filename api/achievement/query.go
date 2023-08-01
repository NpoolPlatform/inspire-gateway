package achievement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	achievement1 "github.com/NpoolPlatform/inspire-gateway/pkg/achievement"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/achievement"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAchievements(ctx context.Context, in *npool.GetAchievementsRequest) (*npool.GetAchievementsResponse, error) {
	handler, err := achievement1.NewHandler(
		ctx,
		achievement1.WithAppID(&in.AppID),
		achievement1.WithUserID(&in.UserID),
		achievement1.WithOffset(in.GetOffset()),
		achievement1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAchievements",
			"In", in,
			"Err", err,
		)
		return &npool.GetAchievementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAchievements(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAchievements",
			"In", in,
			"Err", err,
		)
		return &npool.GetAchievementsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAchievementsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetUserAchievements(ctx context.Context, in *npool.GetUserAchievementsRequest) (*npool.GetUserAchievementsResponse, error) {
	handler, err := achievement1.NewHandler(
		ctx,
		achievement1.WithAppID(&in.AppID),
		achievement1.WithUserIDs(&in.UserIDs),
		achievement1.WithOffset(in.GetOffset()),
		achievement1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAchievements",
			"In", in,
			"Err", err,
		)
		return &npool.GetUserAchievementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAchievements(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAchievements",
			"In", in,
			"Err", err,
		)
		return &npool.GetUserAchievementsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetUserAchievementsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
