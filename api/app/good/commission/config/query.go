//nolint:dupl
package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionconfig1 "github.com/NpoolPlatform/inspire-gateway/pkg/app/good/commission/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/good/commission/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAppGoodCommissionConfigs(ctx context.Context, in *npool.GetAppGoodCommissionConfigsRequest) (*npool.GetAppGoodCommissionConfigsResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.AppID, true),
		commissionconfig1.WithEndAt(in.EndAt, false),
		commissionconfig1.WithOffset(in.GetOffset()),
		commissionconfig1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppGoodCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppGoodCommissionConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCommissions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppGoodCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppGoodCommissionConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppGoodCommissionConfigsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) AdminGetAppGoodCommissionConfigs(ctx context.Context, in *npool.AdminGetAppGoodCommissionConfigsRequest) (*npool.AdminGetAppGoodCommissionConfigsResponse, error) {
	handler, err := commissionconfig1.NewHandler(
		ctx,
		commissionconfig1.WithAppID(&in.TargetAppID, true),
		commissionconfig1.WithEndAt(in.EndAt, false),
		commissionconfig1.WithOffset(in.GetOffset()),
		commissionconfig1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetAppGoodCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetAppGoodCommissionConfigsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCommissions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetAppGoodCommissionConfigs",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetAppGoodCommissionConfigsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetAppGoodCommissionConfigsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
