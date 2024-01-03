package cashcontrol

import (
	"fmt"

	cashcontrolmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/cashcontrol"
	npool "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/app/cashcontrol"

	"context"
)

func (h *Handler) DeleteCashControl(ctx context.Context) (*npool.CashControl, error) {
	info, err := h.GetCashControl(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid cashcontrol")
	}
	if info.ID != *h.ID || info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}
	if _, err := cashcontrolmwcli.DeleteCashControl(ctx, *h.ID); err != nil {
		return nil, err
	}
	return info, nil
}
