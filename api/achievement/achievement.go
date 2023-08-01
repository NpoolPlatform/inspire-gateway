package achievement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/achievement"

	achievement1 "github.com/NpoolPlatform/inspire-gateway/pkg/achievement"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetGoodAchievements(
	ctx context.Context, in *npool.GetGoodAchievementsRequest,
) (
	*npool.GetGoodAchievementsResponse, error,
) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetGoodAchievements", "AppID", in.GetAppID(), "error", err)
		return &npool.GetGoodAchievementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("GetGoodAchievements", "UserID", in.GetUserID(), "error", err)
		return &npool.GetGoodAchievementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := usermwcli.GetUser(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		return &npool.GetGoodAchievementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if user == nil {
		return &npool.GetGoodAchievementsResponse{}, status.Error(codes.InvalidArgument, "User is invalid")
	}

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return &npool.GetGoodAchievementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.GetGoodAchievementsResponse{}, status.Error(codes.InvalidArgument, "App is invalid")
	}

	infos, total, err := achievement1.GetGoodAchievements(ctx,
		in.GetAppID(), in.GetUserID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetGoodAchievements", "AppID", in.GetAppID(), "UserID", in.GetUserID(), "error", err)
		return &npool.GetGoodAchievementsResponse{}, status.Error(codes.Internal, "fail get coin achievements")
	}

	return &npool.GetGoodAchievementsResponse{
		Achievements: infos,
		Total:        total,
	}, nil
}

func (s *Server) GetUserGoodAchievements(
	ctx context.Context, in *npool.GetUserGoodAchievementsRequest,
) (
	*npool.GetUserGoodAchievementsResponse, error,
) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetUserGoodAchievements", "AppID", in.GetAppID(), "error", err)
		return &npool.GetUserGoodAchievementsResponse{}, status.Error(codes.Internal, "AppID is invalid")
	}

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return &npool.GetUserGoodAchievementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.GetUserGoodAchievementsResponse{}, status.Error(codes.InvalidArgument, "App is invalid")
	}

	if len(in.GetUserIDs()) == 0 {
		logger.Sugar().Errorw("GetUserGoodAchievements", "UserIDs", in.GetUserIDs(), "error", "UserIDs is invalid")
		return &npool.GetUserGoodAchievementsResponse{}, status.Error(codes.Internal, "UserIDs is invalid")
	}

	for _, user := range in.GetUserIDs() {
		if _, err := uuid.Parse(user); err != nil {
			logger.Sugar().Errorw("GetUserGoodAchievements", "UserID", user, "error", err)
			return &npool.GetUserGoodAchievementsResponse{}, status.Error(codes.Internal, "UserIDs is invalid")
		}
	}

	infos, total, err := achievement1.GetUserGoodAchievements(ctx,
		in.GetAppID(), in.GetUserIDs(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetUserGoodAchievements", "AppID", in.GetAppID(), "UserIDs", in.GetUserIDs(), "error", err)
		return &npool.GetUserGoodAchievementsResponse{}, status.Error(codes.Internal, "fail get coin achievements")
	}

	return &npool.GetUserGoodAchievementsResponse{
		Achievements: infos,
		Total:        total,
	}, nil
}
