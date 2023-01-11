package allocated

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/allocated"
	allocatedmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"
	allocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/allocated"

	allocated1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/allocated"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/coupon"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

//nolint
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

	user, err := usermwcli.GetUser(ctx, in.GetAppID(), in.GetTargetUserID())
	if err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if user == nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "TargetUserID is invalid")
	}

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	coup, err := couponmwcli.GetCoupon(ctx, in.GetCouponID(), in.GetCouponType())
	if err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if coup == nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "CouponID is invalid")
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
