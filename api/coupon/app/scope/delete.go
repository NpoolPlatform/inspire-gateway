package scope

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	appgoodscope1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/scope"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteAppGoodScope(ctx context.Context, in *npool.DeleteAppGoodScopeRequest) (*npool.DeleteAppGoodScopeResponse, error) {
	handler, err := appgoodscope1.NewHandler(
		ctx,
		appgoodscope1.WithID(&in.ID, true),
		appgoodscope1.WithEntID(&in.EntID, true),
		appgoodscope1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppGoodScope",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppGoodScopeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteAppGoodScope(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppGoodScope",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppGoodScopeResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAppGoodScopeResponse{
		Info: info,
	}, nil
}
