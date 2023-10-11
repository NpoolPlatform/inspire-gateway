package scope

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	scope1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteScope(ctx context.Context, in *npool.DeleteScopeRequest) (*npool.DeleteScopeResponse, error) {
	handler, err := scope1.NewHandler(
		ctx,
		scope1.WithID(&in.ID, true),
		scope1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteScope",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteScopeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteScope(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteScope",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteScopeResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteScopeResponse{
		Info: info,
	}, nil
}
