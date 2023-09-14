//nolint:dupl
package commission

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commission1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CloneCommissions(ctx context.Context, in *npool.CloneCommissionsRequest) (*npool.CloneCommissionsResponse, error) {
	handler, err := commission1.NewHandler(
		ctx,
		commission1.WithAppID(&in.AppID),
		commission1.WithFromAppGoodID(&in.FromAppGoodID),
		commission1.WithToAppGoodID(&in.ToAppGoodID),
		commission1.WithScalePercent(&in.ScalePercent),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CloneCommissions",
			"In", in,
			"Err", err,
		)
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = handler.CloneCommissions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CloneCommissions",
			"In", in,
			"Err", err,
		)
		return &npool.CloneCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CloneCommissionsResponse{}, nil
}

func (s *Server) CloneAppCommissions(ctx context.Context, in *npool.CloneAppCommissionsRequest) (*npool.CloneAppCommissionsResponse, error) {
	handler, err := commission1.NewHandler(
		ctx,
		commission1.WithAppID(&in.TargetAppID),
		commission1.WithFromAppGoodID(&in.FromAppGoodID),
		commission1.WithToAppGoodID(&in.ToAppGoodID),
		commission1.WithScalePercent(&in.ScalePercent),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CloneCommissions",
			"In", in,
			"Err", err,
		)
		return &npool.CloneAppCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = handler.CloneCommissions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CloneCommissions",
			"In", in,
			"Err", err,
		)
		return &npool.CloneAppCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CloneAppCommissionsResponse{}, nil
}
