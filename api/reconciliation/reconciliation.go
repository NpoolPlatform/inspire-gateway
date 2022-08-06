package reconciliation

import (
	"context"

	reconciliation1 "github.com/NpoolPlatform/inspire-gateway/pkg/reconciliation"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/reconciliation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) UpdateArchivement(ctx context.Context, in *npool.UpdateArchivementRequest) (*npool.UpdateArchivementResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.UpdateArchivementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		return &npool.UpdateArchivementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := reconciliation1.UpdateArchivement(ctx, in.GetAppID(), in.GetUserID()); err != nil {
		return &npool.UpdateArchivementResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateArchivementResponse{}, nil
}
