package allocated

import (
	"context"

	allocatedmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/allocated"
	allocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/allocated"
)

func (h *Handler) CreateCoupon(ctx context.Context) (*allocatedmwpb.Coupon, error) {
	return allocatedmwcli.CreateCoupon(
		ctx,
		&allocatedmwpb.CouponReq{
			ID:       h.ID,
			AppID:    h.AppID,
			UserID:   h.UserID,
			CouponID: h.CouponID,
		},
	)
}
