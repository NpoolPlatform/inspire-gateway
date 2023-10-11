package scope

import (
	"fmt"

	scopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"

	"context"
)

func (h *Handler) DeleteScope(ctx context.Context) (*npool.Scope, error) {
	info, err := h.GetScope(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	if info.AppID != *h.AppID {
		return nil, fmt.Errorf("app id not matched")
	}
	if _, err := scopemwcli.DeleteScope(ctx, *h.ID); err != nil {
		return nil, err
	}
	return info, nil
}
