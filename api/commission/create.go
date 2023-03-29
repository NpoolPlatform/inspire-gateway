package commission

import (
	"context"
	"fmt"
	"time"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	commmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"

	comm1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"

	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	regmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	"github.com/shopspring/decimal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) createCommission(ctx context.Context, in *npool.CreateCommissionRequest) (*npool.CreateCommissionResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	switch in.GetSettleType() {
	case commmgrpb.SettleType_GoodOrderPercent:
		fallthrough //nolint
	case commmgrpb.SettleType_GoodOrderValuePercent:
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

	startAt := in.GetStartAt()
	if startAt == 0 {
		startAt = uint32(time.Now().Unix())
	}

	info, err := comm1.CreateCommission(
		ctx,
		in.GetAppID(),
		in.GetTargetUserID(),
		in.GoodID,
		in.GetSettleType(),
		value,
		&startAt,
	)
	if err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCommissionResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateCommission(ctx context.Context, in *npool.CreateCommissionRequest) (*npool.CreateCommissionResponse, error) {
	conds := &commmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		UserID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetUserID(),
		},
		SettleType: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(in.GetSettleType()),
		},
		EndAt: &commonpb.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(0),
		},
	}
	if in.GoodID != nil {
		conds.GoodID = &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetGoodID(),
		}
	}
	// Only users with the same type of commission can set commissions for subordinate users
	comm, err := commmwcli.GetCommissionOnly(ctx, conds)
	if err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}
	if comm == nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.Internal, "commission not exist")
	}

	value, err := decimal.NewFromString(in.GetValue())
	if err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	switch in.GetSettleType() {
	case commmgrpb.SettleType_GoodOrderPercent:
		fallthrough //nolint
	case commmgrpb.SettleType_GoodOrderValuePercent:
		percent, err := decimal.NewFromString(comm.GetPercent())
		if err != nil {
			return nil, err
		}
		if percent.Cmp(value) < 0 {
			return nil, fmt.Errorf("invalid percent")
		}
	case commmgrpb.SettleType_LimitedOrderPercent:
		fallthrough //nolint
	case commmgrpb.SettleType_AmountThreshold:
		fallthrough //nolint
	case commmgrpb.SettleType_NoCommission:
		return nil, fmt.Errorf("not implemented")
	default:
		return nil, fmt.Errorf("unknown settle type")
	}

	reg, err := regmwcli.GetRegistrationOnly(ctx, &regmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		InviterID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetUserID(),
		},
		InviteeID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetUserID(),
		},
	})
	if err != nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}
	if reg == nil {
		return &npool.CreateCommissionResponse{}, status.Error(codes.Internal, "permission denied")
	}

	return s.createCommission(ctx, in)
}

func (s *Server) CreateUserCommission(
	ctx context.Context,
	in *npool.CreateUserCommissionRequest,
) (
	*npool.CreateUserCommissionResponse,
	error,
) {
	resp, err := s.createCommission(ctx, &npool.CreateCommissionRequest{
		AppID:        in.AppID,
		TargetUserID: in.TargetUserID,
		GoodID:       in.GoodID,
		SettleType:   in.SettleType,
		Value:        in.Value,
		StartAt:      in.StartAt,
	})
	if err != nil {
		return &npool.CreateUserCommissionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateUserCommissionResponse{
		Info: resp.Info,
	}, nil
}
