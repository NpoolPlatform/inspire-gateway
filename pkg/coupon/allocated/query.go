package allocated

import (
	"context"

	allocatedmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/allocated"
	allocatedmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"
	allocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/allocated"
)

func GetCoupons(ctx context.Context, conds *allocatedmgrpb.Conds, offset, limit int32) ([]*allocatedmwpb.Coupon, uint32, error) {
	return allocatedmwcli.GetCoupons(ctx, conds, offset, limit)
}
