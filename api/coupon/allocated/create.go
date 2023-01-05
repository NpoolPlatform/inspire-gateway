package allocated

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/allocated"
	allocatedmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"
	allocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/allocated"

	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/allocated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) CreateCoupon(ctx context.Context, in *npool.CreateCouponRequest) (*npool.CreateCouponResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetCouponID()); err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	switch in.GetCouponType() {
	case allocatedmgrpb.CouponType_FixAmount:
	case allocatedmgrpb.CouponType_Discount:
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
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "Not implemented")
	default:
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "Unknown coupon type")
	}

	info, err := allocated1.CreateCoupon(ctx, &allocatedmwpb.CouponReq{
		CouponType: &in.CouponType,
		AppID:      &in.AppID,
		CouponID:   &in.CouponID,
		UserID:     &in.TargetUserID,
	})
	if err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCouponResponse{
		Info: info,
	}, nil
}
