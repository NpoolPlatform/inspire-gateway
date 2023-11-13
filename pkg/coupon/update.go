package coupon

import (
	"context"
	"fmt"

	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"
)

func (h *Handler) UpdateCoupon(ctx context.Context) (*couponmwpb.Coupon, error) {
	info, err := h.GetCoupon(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	if info.AppID != *h.AppID || info.EntID != *h.EntID {
		return nil, fmt.Errorf("permission denied")
	}

	return couponmwcli.UpdateCoupon(ctx, &couponmwpb.CouponReq{
		ID:               h.ID,
		Denomination:     h.Denomination,
		Circulation:      h.Circulation,
		IssuedBy:         h.IssuedBy,
		StartAt:          h.StartAt,
		DurationDays:     h.DurationDays,
		Message:          h.Message,
		Name:             h.Name,
		CouponConstraint: h.CouponConstraint,
		Threshold:        h.Threshold,
		Random:           h.Random,
		CouponScope:      h.CouponScope,
	})
}
