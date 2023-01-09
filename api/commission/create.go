package commission

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	commmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"

	comm1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"

	"github.com/shopspring/decimal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) CreateCommission(ctx context.Context, in *npool.CreateCommissionRequest) (*npool.CreateCommissionResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	switch in.GetSettleType() {
	case commmgrpb.SettleType_GoodOrderPercent:
		if in.GoodID != nil {
			if _, err := uuid.Parse(in.GetGoodID()); err != nil {
				return &npool.CreateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
			}
		}
	case commmgrpb.SettleType_LimitedOrderPercent:
		fallthrough //nolint
	case commmgrpb.SettleType_AmountThreshold:
		fallthrough //nolint
	case commmgrpb.SettleType_NoCommission:
		return &npool.CreateCommissionResponse{}, status.Error(codes.InvalidArgument, "NOT IMPLEMENTED")
	default:
		return &npool.CreateCommissionResponse{}, status.Error(codes.InvalidArgument, "Unknown settle type")
	}

	value, err := decimal.NewFromString(in.GetValue())
	if err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := comm1.CreateCommission(
		ctx,
		in.GetAppID(),
		in.GetTargetUserID(),
		in.GoodID,
		in.GetSettleType(),
		value,
		in.StartAt,
	)
	if err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCommissionResponse{
		Info: info,
	}, nil
}
