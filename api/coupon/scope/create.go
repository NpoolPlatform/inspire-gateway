package scope

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	scope1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateScope(ctx context.Context, in *npool.CreateScopeRequest) (*npool.CreateScopeResponse, error) {
	handler, err := scope1.NewHandler(
		ctx,
		scope1.WithGoodID(&in.GoodID, true),
		scope1.WithCouponID(&in.CouponID, true),
		scope1.WithCouponScope(&in.CouponScope, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateScope",
			"In", in,
			"Err", err,
		)
		return &npool.CreateScopeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateScope(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateScope",
			"In", in,
			"Err", err,
		)
		return &npool.CreateScopeResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateScopeResponse{
		Info: info,
	}, nil
}
