package allocated

import (
	"context"

	allocatedmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/allocated"
	allocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/allocated"
)

func CreateCoupon(ctx context.Context, in *allocatedmwpb.CouponReq) (*allocatedmwpb.Coupon, error) {
	return allocatedmwcli.CreateCoupon(ctx, in)
}
