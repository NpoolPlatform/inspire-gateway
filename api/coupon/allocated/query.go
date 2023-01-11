package allocated

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/allocated"
	allocatedmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/allocated"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetCoupons(ctx context.Context, in *npool.GetCouponsRequest) (*npool.GetCouponsResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.GetCouponsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		return &npool.GetCouponsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := constant.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := allocated1.GetCoupons(ctx, &allocatedmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		UserID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetUserID(),
		},
	}, in.GetOffset(), limit)
	if err != nil {
		return &npool.GetCouponsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCouponsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppCoupons(ctx context.Context, in *npool.GetAppCouponsRequest) (*npool.GetAppCouponsResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.GetAppCouponsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := constant.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := allocated1.GetCoupons(ctx, &allocatedmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
	}, in.GetOffset(), limit)
	if err != nil {
		return &npool.GetAppCouponsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCouponsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
