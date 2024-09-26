package reward

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/user/reward"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/user/reward"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetMyUserRewards(ctx context.Context, in *npool.GetMyUserRewardsRequest) (*npool.GetMyUserRewardsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.AppID, true),
		allocated1.WithUserID(&in.UserID, false),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMyUserRewards",
			"In", in,
			"Err", err,
		)
		return &npool.GetMyUserRewardsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetUserRewards(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMyUserRewards",
			"In", in,
			"Err", err,
		)
		return &npool.GetMyUserRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetMyUserRewardsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) AdminGetUserRewards(ctx context.Context, in *npool.AdminGetUserRewardsRequest) (*npool.AdminGetUserRewardsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.TargetAppID, true),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetUserRewards",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetUserRewardsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetUserRewards(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetUserRewards",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetUserRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetUserRewardsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
