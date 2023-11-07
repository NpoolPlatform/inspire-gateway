package scope

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	appgoodscope1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/app/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/scope"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateAppGoodScope(ctx context.Context, in *npool.CreateAppGoodScopeRequest) (*npool.CreateAppGoodScopeResponse, error) {
	handler, err := appgoodscope1.NewHandler(
		ctx,
		appgoodscope1.WithAppID(&in.AppID, true),
		appgoodscope1.WithAppGoodID(&in.AppGoodID, true),
		appgoodscope1.WithScopeID(&in.ScopeID, true),
		appgoodscope1.WithCouponScope(in.CouponScope, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppGoodScope",
			"In", in,
			"Err", err,
		)
		return &npool.CreateAppGoodScopeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateAppGoodScope(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppGoodScope",
			"In", in,
			"Err", err,
		)
		return &npool.CreateAppGoodScopeResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppGoodScopeResponse{
		Info: info,
	}, nil
}
