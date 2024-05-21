package allocated

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/coin/allocated"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/allocated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UserGetCoinAllocateds(ctx context.Context, in *npool.UserGetCoinAllocatedsRequest) (*npool.UserGetCoinAllocatedsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.AppID, true),
		allocated1.WithUserID(&in.UserID, false),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UserGetCoinAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.UserGetCoinAllocatedsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoinAllocateds(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UserGetCoinAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.UserGetCoinAllocatedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UserGetCoinAllocatedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) AdminGetAppCoinAllocateds(ctx context.Context, in *npool.AdminGetAppCoinAllocatedsRequest) (*npool.AdminGetAppCoinAllocatedsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.TargetAppID, true),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetAppCoinAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetAppCoinAllocatedsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoinAllocateds(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetAppCoinAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetAppCoinAllocatedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetAppCoinAllocatedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
