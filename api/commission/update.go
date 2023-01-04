package commission

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	comm1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"

	"github.com/shopspring/decimal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) UpdateCommission(ctx context.Context, in *npool.UpdateCommissionRequest) (*npool.UpdateCommissionResponse, error) {
	if _, err := uuid.Parse(in.GetID()); err != nil {
		return &npool.UpdateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.UpdateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.Value == nil && in.StartAt == nil {
		return &npool.UpdateCommissionResponse{}, status.Error(codes.InvalidArgument, "Nothing to be done")
	}

	if in.Value != nil {
		if _, err := decimal.NewFromString(in.GetValue()); err != nil {
			return &npool.UpdateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	info, err := comm1.UpdateCommission(
		ctx,
		in.GetID(),
		in.GetAppID(),
		in.Value,
		in.StartAt,
	)
	if err != nil {
		return &npool.UpdateCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCommissionResponse{
		Info: info,
	}, nil
}
