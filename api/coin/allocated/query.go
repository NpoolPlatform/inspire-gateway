package allocated

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/coin/allocated"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/allocated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetMyCoinAllocateds(ctx context.Context, in *npool.GetMyCoinAllocatedsRequest) (*npool.GetMyCoinAllocatedsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.AppID, true),
		allocated1.WithUserID(&in.UserID, false),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMyCoinAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.GetMyCoinAllocatedsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoinAllocateds(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMyCoinAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.GetMyCoinAllocatedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetMyCoinAllocatedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) AdminGetCoinAllocateds(ctx context.Context, in *npool.AdminGetCoinAllocatedsRequest) (*npool.AdminGetCoinAllocatedsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.TargetAppID, true),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetCoinAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetCoinAllocatedsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoinAllocateds(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetCoinAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetCoinAllocatedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetCoinAllocatedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
