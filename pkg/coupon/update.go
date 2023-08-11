package coupon

import (
	"context"
	"fmt"

	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"
)

func (h *Handler) UpdateCoupon(ctx context.Context) (*couponmwpb.Coupon, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
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
		GoodID:           h.GoodID,
		CouponConstraint: h.CouponConstraint,
		Threshold:        h.Threshold,
		Random:           h.Random,
	})
}
