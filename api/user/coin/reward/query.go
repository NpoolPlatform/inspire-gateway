package reward

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/user/coin/reward"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/user/coin/reward"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UserGetUserCoinRewards(ctx context.Context, in *npool.UserGetUserCoinRewardsRequest) (*npool.UserGetUserCoinRewardsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.AppID, true),
		allocated1.WithUserID(&in.UserID, false),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UserGetUserCoinRewards",
			"In", in,
			"Err", err,
		)
		return &npool.UserGetUserCoinRewardsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetUserCoinRewards(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UserGetUserCoinRewards",
			"In", in,
			"Err", err,
		)
		return &npool.UserGetUserCoinRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UserGetUserCoinRewardsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) AdminGetUserCoinRewards(ctx context.Context, in *npool.AdminGetUserCoinRewardsRequest) (*npool.AdminGetUserCoinRewardsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.TargetAppID, true),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetUserCoinRewards",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetUserCoinRewardsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetUserCoinRewards(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetUserCoinRewards",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetUserCoinRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetUserCoinRewardsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
