package scope

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	appgoodscope1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/scope"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAppGoodScopes(ctx context.Context, in *npool.GetAppGoodScopesRequest) (*npool.GetAppGoodScopesResponse, error) {
	handler, err := appgoodscope1.NewHandler(
		ctx,
		appgoodscope1.WithAppID(&in.AppID, true),
		appgoodscope1.WithOffset(in.GetOffset()),
		appgoodscope1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppGoodScopes",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppGoodScopesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAppGoodScopes(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppGoodScopes",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppGoodScopesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppGoodScopesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
