package reconcile

import (
	"context"

	reconcile1 "github.com/NpoolPlatform/inspire-gateway/pkg/reconcile"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/reconcile"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) Reconcile(ctx context.Context, in *npool.ReconcileRequest) (*npool.ReconcileResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.ReconcileResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		return &npool.ReconcileResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := uuid.Parse(in.GetGoodID()); err != nil {
		return &npool.ReconcileResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := reconcile1.Reconcile(ctx, in.GetAppID(), in.GetUserID(), in.GetGoodID()); err != nil {
		return &npool.ReconcileResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ReconcileResponse{}, nil
}
