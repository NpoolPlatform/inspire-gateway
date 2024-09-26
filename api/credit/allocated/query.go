package allocated

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/credit/allocated"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/credit/allocated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetMyCreditAllocateds(ctx context.Context, in *npool.GetMyCreditAllocatedsRequest) (*npool.GetMyCreditAllocatedsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.AppID, true),
		allocated1.WithUserID(&in.UserID, false),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMyCreditAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.GetMyCreditAllocatedsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCreditAllocateds(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetMyCreditAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.GetMyCreditAllocatedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetMyCreditAllocatedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) AdminGetCreditAllocateds(ctx context.Context, in *npool.AdminGetCreditAllocatedsRequest) (*npool.AdminGetCreditAllocatedsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.TargetAppID, true),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetCreditAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetCreditAllocatedsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCreditAllocateds(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetCreditAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetCreditAllocatedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetCreditAllocatedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
