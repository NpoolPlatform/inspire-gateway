package coin

import (
	"fmt"

	couponcoinmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/coin"

	"context"
)

func (h *Handler) DeleteCouponCoin(ctx context.Context) (*npool.CouponCoin, error) {
	info, err := h.GetCouponCoin(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid couponcoin")
	}
	if info.ID != *h.ID || info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}
	if _, err := couponcoinmwcli.DeleteCouponCoin(ctx, *h.ID); err != nil {
		return nil, err
	}
	return info, nil
}
