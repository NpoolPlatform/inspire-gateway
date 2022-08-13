package archivement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/archivement"

	archivement1 "github.com/NpoolPlatform/inspire-gateway/pkg/archivement"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetGoodArchivements(
	ctx context.Context, in *npool.GetGoodArchivementsRequest,
) (
	*npool.GetGoodArchivementsResponse, error,
) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetGoodArchivements", "AppID", in.GetAppID(), "error", err)
		return &npool.GetGoodArchivementsResponse{}, status.Error(codes.Internal, "AppID is invalid")
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("GetGoodArchivements", "UserID", in.GetUserID(), "error", err)
		return &npool.GetGoodArchivementsResponse{}, status.Error(codes.Internal, "UserID is invalid")
	}

	infos, total, err := archivement1.GetGoodArchivements(ctx,
		in.GetAppID(), in.GetUserID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetGoodArchivements", "AppID", in.GetAppID(), "UserID", in.GetUserID(), "error", err)
		return &npool.GetGoodArchivementsResponse{}, status.Error(codes.Internal, "fail get coin archivements")
	}

	return &npool.GetGoodArchivementsResponse{
		Archivements: infos,
		Total:        total,
	}, nil
}

func (s *Server) GetUserGoodArchivements(
	ctx context.Context, in *npool.GetUserGoodArchivementsRequest,
) (
	*npool.GetUserGoodArchivementsResponse, error,
) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetUserGoodArchivements", "AppID", in.GetAppID(), "error", err)
		return &npool.GetUserGoodArchivementsResponse{}, status.Error(codes.Internal, "AppID is invalid")
	}

	if len(in.GetUserIDs()) == 0 {
		logger.Sugar().Errorw("GetUserGoodArchivements", "UserIDs", in.GetUserIDs(), "error", "UserIDs is invalid")
		return &npool.GetUserGoodArchivementsResponse{}, status.Error(codes.Internal, "UserIDs is invalid")
	}

	for _, user := range in.GetUserIDs() {
		if _, err := uuid.Parse(user); err != nil {
			logger.Sugar().Errorw("GetUserGoodArchivements", "UserID", user, "error", err)
			return &npool.GetUserGoodArchivementsResponse{}, status.Error(codes.Internal, "UserIDs is invalid")
		}
	}

	infos, total, err := archivement1.GetUserGoodArchivements(ctx,
		in.GetAppID(), in.GetUserIDs(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetUserGoodArchivements", "AppID", in.GetAppID(), "UserIDs", in.GetUserIDs(), "error", err)
		return &npool.GetUserGoodArchivementsResponse{}, status.Error(codes.Internal, "fail get coin archivements")
	}

	return &npool.GetUserGoodArchivementsResponse{
		Archivements: infos,
		Total:        total,
	}, nil
}
