package allocated

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/credit/allocated"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/credit/allocated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UserGetCreditAllocateds(ctx context.Context, in *npool.UserGetCreditAllocatedsRequest) (*npool.UserGetCreditAllocatedsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.AppID, true),
		allocated1.WithUserID(&in.UserID, false),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UserGetCreditAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.UserGetCreditAllocatedsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCreditAllocateds(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UserGetCreditAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.UserGetCreditAllocatedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UserGetCreditAllocatedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) AdminGetAppCreditAllocateds(ctx context.Context, in *npool.AdminGetAppCreditAllocatedsRequest) (*npool.AdminGetAppCreditAllocatedsResponse, error) {
	handler, err := allocated1.NewHandler(
		ctx,
		allocated1.WithAppID(&in.TargetAppID, true),
		allocated1.WithOffset(in.GetOffset()),
		allocated1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetAppCreditAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetAppCreditAllocatedsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCreditAllocateds(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetAppCreditAllocateds",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetAppCreditAllocatedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetAppCreditAllocatedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
