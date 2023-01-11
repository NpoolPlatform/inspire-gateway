package coupon

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/coupon"
	allocatedmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/coupon"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	coupon1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/coupon"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetCoupons(ctx context.Context, in *npool.GetCouponsRequest) (*npool.GetCouponsResponse, error) {
	switch in.GetCouponType() {
	case allocatedmgrpb.CouponType_FixAmount:
		fallthrough //nolint
	case allocatedmgrpb.CouponType_Discount:
		fallthrough //nolint
	case allocatedmgrpb.CouponType_SpecialOffer:
	case allocatedmgrpb.CouponType_ThresholdFixAmount:
		fallthrough //nolint
	case allocatedmgrpb.CouponType_ThresholdDiscount:
		fallthrough //nolint
	case allocatedmgrpb.CouponType_GoodFixAmount:
		fallthrough //nolint
	case allocatedmgrpb.CouponType_GoodDiscount:
		fallthrough //nolint
	case allocatedmgrpb.CouponType_GoodThresholdFixAmount:
		fallthrough //nolint
	case allocatedmgrpb.CouponType_GoodThresholdDiscount:
		return &npool.GetCouponsResponse{}, status.Error(codes.InvalidArgument, "Not implemented")
	default:
		return &npool.GetCouponsResponse{}, status.Error(codes.InvalidArgument, "Unknown coupon type")
	}

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.GetCouponsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := constant.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := coupon1.GetCoupons(ctx, &couponmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		CouponType: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(in.GetCouponType()),
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
	resp, err := s.GetCoupons(ctx, &npool.GetCouponsRequest{
		AppID:      in.GetTargetAppID(),
		CouponType: in.GetCouponType(),
		Offset:     in.GetOffset(),
		Limit:      in.GetLimit(),
	})
	if err != nil {
		return &npool.GetAppCouponsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCouponsResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}
