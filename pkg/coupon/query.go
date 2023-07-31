package coupon

import (
	"context"

	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/coupon"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/coupon"
)

func GetCoupons(ctx context.Context, conds *couponmwpb.Conds, offset, limit int32) ([]*couponmwpb.Coupon, uint32, error) {
	return couponmwcli.GetCoupons(ctx, conds, offset, limit)
}
