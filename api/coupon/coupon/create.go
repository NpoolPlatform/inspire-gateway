package coupon

import (
	"context"
	"time"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/coupon"
	allocatedmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/coupon"

	coupon1 "github.com/NpoolPlatform/inspire-gateway/pkg/coupon/coupon"

	timedef "github.com/NpoolPlatform/go-service-framework/pkg/const/time"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	"github.com/shopspring/decimal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

//nolint
func (s *Server) CreateCoupon(ctx context.Context, in *npool.CreateCouponRequest) (*npool.CreateCouponResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	switch in.GetCouponType() {
	case allocatedmgrpb.CouponType_FixAmount:
		fallthrough //nolint
	case allocatedmgrpb.CouponType_Discount:
	case allocatedmgrpb.CouponType_SpecialOffer:
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "Not supported")
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

	if _, err := decimal.NewFromString(in.GetValue()); err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := decimal.NewFromString(in.GetCirculation()); err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.GetStartAt()+in.GetDurationDays()*timedef.SecondsPerDay <= uint32(time.Now().Unix()) {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "StartAt is invalid")
	}
	if in.GetDurationDays() == 0 {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "DurationDays is invalid")
	}
	if in.GetMessage() == "" {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "Message is invalid")
	}
	if in.GetName() == "" {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "Name is invalid")
	}

	switch in.GetCouponType() {
	case allocatedmgrpb.CouponType_FixAmount:
	case allocatedmgrpb.CouponType_Discount:
	case allocatedmgrpb.CouponType_SpecialOffer:
		if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
			return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		user, err := usermwcli.GetUser(ctx, in.GetTargetAppID(), in.GetTargetUserID())
		if err != nil {
			return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		if user == nil {
			return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "TargetUserID is invalid")
		}
	case allocatedmgrpb.CouponType_ThresholdFixAmount:
	case allocatedmgrpb.CouponType_ThresholdDiscount:
	case allocatedmgrpb.CouponType_GoodFixAmount:
	case allocatedmgrpb.CouponType_GoodDiscount:
	case allocatedmgrpb.CouponType_GoodThresholdFixAmount:
	case allocatedmgrpb.CouponType_GoodThresholdDiscount:
	default:
	}

	user, err := usermwcli.GetUser(ctx, in.GetAppID(), in.GetUserID())
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

	app, err = appmwcli.GetApp(ctx, in.GetTargetAppID())
	if err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.InvalidArgument, "TargetAppID is invalid")
	}

	req := &couponmwpb.CouponReq{
		CouponType:       &in.CouponType,
		AppID:            &in.TargetAppID,
		Value:            &in.Value,
		Circulation:      &in.Circulation,
		ReleasedByUserID: &in.UserID,
		StartAt:          &in.StartAt,
		DurationDays:     &in.DurationDays,
		Message:          &in.Message,
		Name:             &in.Name,
	}

	switch in.GetCouponType() {
	case allocatedmgrpb.CouponType_FixAmount:
	case allocatedmgrpb.CouponType_Discount:
	case allocatedmgrpb.CouponType_SpecialOffer:
		req.UserID = in.TargetUserID
	case allocatedmgrpb.CouponType_ThresholdFixAmount:
	case allocatedmgrpb.CouponType_ThresholdDiscount:
	case allocatedmgrpb.CouponType_GoodFixAmount:
	case allocatedmgrpb.CouponType_GoodDiscount:
	case allocatedmgrpb.CouponType_GoodThresholdFixAmount:
	case allocatedmgrpb.CouponType_GoodThresholdDiscount:
	default:
	}

	info, err := coupon1.CreateCoupon(ctx, req)
	if err != nil {
		return &npool.CreateCouponResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCouponResponse{
		Info: info,
	}, nil
}
