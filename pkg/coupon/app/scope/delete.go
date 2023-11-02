package scope

import (
	"fmt"

	appgoodscopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/scope"

	"context"
)

func (h *Handler) DeleteAppGoodScope(ctx context.Context) (*npool.Scope, error) {
	info, err := h.GetAppGoodScope(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid scope")
	}
	if _, err := appgoodscopemwcli.DeleteAppGoodScope(ctx, *h.ID); err != nil {
		return nil, err
	}
	return info, nil
}
