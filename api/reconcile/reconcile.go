package reconcile

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	reconcile1 "github.com/NpoolPlatform/inspire-gateway/pkg/reconcile"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/reconcile"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Reconcile(ctx context.Context, in *npool.ReconcileRequest) (*npool.ReconcileResponse, error) {
	handler, err := reconcile1.NewHandler(
		ctx,
		reconcile1.WithAppID(&in.AppID),
		reconcile1.WithAppGoodID(&in.AppGoodID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"Reconcile",
			"In", in,
			"Err", err,
		)
		return &npool.ReconcileResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := handler.Reconcile(ctx); err != nil {
		logger.Sugar().Errorw(
			"Reconcile",
			"In", in,
			"Err", err,
		)
		return &npool.ReconcileResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ReconcileResponse{}, nil
}
