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
	return couponmwcli.GetCoupons(
		ctx,
		&couponmwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		},
		h.Offset,
		h.Limit,
	)
}
