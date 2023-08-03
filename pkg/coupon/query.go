package coupon

import (
	"context"
	"fmt"

	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"
)

func (h *Handler) GetCoupons(ctx context.Context) ([]*couponmwpb.Coupon, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appid")
	}
	conds := &couponmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.CouponType != nil {
		conds.CouponType = &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.CouponType)}
	}
	return couponmwcli.GetCoupons(ctx, conds, h.Offset, h.Limit)
}
