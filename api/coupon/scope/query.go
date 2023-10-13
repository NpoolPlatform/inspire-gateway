package scope

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	scope1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetScopes(ctx context.Context, in *npool.GetScopesRequest) (*npool.GetScopesResponse, error) {
	handler, err := scope1.NewHandler(
		ctx,
		scope1.WithAppID(&in.AppID, true),
		scope1.WithOffset(in.GetOffset()),
		scope1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetScopes",
			"In", in,
			"Err", err,
		)
		return &npool.GetScopesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetScopes(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetScopes",
			"In", in,
			"Err", err,
		)
		return &npool.GetScopesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetScopesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
