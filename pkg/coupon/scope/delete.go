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
		return nil, fmt.Errorf("invalid scope")
	}
	if info.ID != *h.ID {
		return nil, fmt.Errorf("permission denied")
	}
	if _, err := scopemwcli.DeleteScope(ctx, *h.ID); err != nil {
		return nil, err
	}
	return info, nil
}
