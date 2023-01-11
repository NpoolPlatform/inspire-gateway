package coupon

import (
	"context"

	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/coupon"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/coupon"
)

func CreateCoupon(ctx context.Context, in *couponmwpb.CouponReq) (*couponmwpb.Coupon, error) {
	return couponmwcli.CreateCoupon(ctx, in)
}
