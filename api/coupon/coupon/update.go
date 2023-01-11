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
	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/coupon"

	"github.com/shopspring/decimal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

//nolint
func (s *Server) UpdateCoupon(ctx context.Context, in *npool.UpdateCouponRequest) (*npool.UpdateCouponResponse, error) {
	if _, err := uuid.Parse(in.GetID()); err != nil {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

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
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "Not implemented")
	default:
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "Unknown coupon type")
	}

	info, err := couponmwcli.GetCoupon(ctx, in.GetID(), in.GetCouponType())
	if err != nil {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if info == nil {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "Coupon is invalid")
	}

	allocated, err := decimal.NewFromString(info.Allocated)
	if err != nil {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if allocated.Cmp(decimal.NewFromInt(0)) > 0 && (in.Value != nil || in.Circulation != nil) {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "permission denied")
	}

	if in.Value != nil {
		if _, err := decimal.NewFromString(in.GetValue()); err != nil {
			return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	if in.Circulation != nil {
		if _, err := decimal.NewFromString(in.GetCirculation()); err != nil {
			return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	startAt := info.StartAt
	if in.StartAt != nil {
		startAt = in.GetStartAt()
	}
	durationDays := info.DurationDays
	if in.DurationDays != nil {
		durationDays = in.GetDurationDays()
	}

	if in.StartAt != nil || in.DurationDays != nil {
		if startAt+durationDays*timedef.SecondsPerDay <= uint32(time.Now().Unix()) {
			return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "StartAt is invalid")
		}
	}
	if in.DurationDays != nil && in.GetDurationDays() == 0 {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "DurationDays is invalid")
	}

	if in.Message != nil && in.GetMessage() == "" {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "Message is invalid")
	}
	if in.Name != nil && in.GetName() == "" {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "Name is invalid")
	}

	switch in.GetCouponType() {
	case allocatedmgrpb.CouponType_FixAmount:
	case allocatedmgrpb.CouponType_Discount:
	case allocatedmgrpb.CouponType_SpecialOffer:
		if in.TargetUserID != nil {
			if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
				return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
			}
			user, err := usermwcli.GetUser(ctx, in.GetTargetAppID(), in.GetTargetUserID())
			if err != nil {
				return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
			}
			if user == nil {
				return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "TargetUserID is invalid")
			}
		}
	case allocatedmgrpb.CouponType_ThresholdFixAmount:
	case allocatedmgrpb.CouponType_ThresholdDiscount:
	case allocatedmgrpb.CouponType_GoodFixAmount:
	case allocatedmgrpb.CouponType_GoodDiscount:
	case allocatedmgrpb.CouponType_GoodThresholdFixAmount:
	case allocatedmgrpb.CouponType_GoodThresholdDiscount:
	default:
	}

	app, err := appmwcli.GetApp(ctx, in.GetTargetAppID())
	if err != nil {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.UpdateCouponResponse{}, status.Error(codes.InvalidArgument, "TargetAppID is invalid")
	}

	req := &couponmwpb.CouponReq{
		ID:           &in.ID,
		AppID:        &in.TargetAppID,
		CouponType:   &in.CouponType,
		Value:        in.Value,
		Circulation:  in.Circulation,
		StartAt:      in.StartAt,
		DurationDays: in.DurationDays,
		Message:      in.Message,
		Name:         in.Name,
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

	info, err = coupon1.UpdateCoupon(ctx, req)
	if err != nil {
		return &npool.UpdateCouponResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCouponResponse{
		Info: info,
	}, nil
}
