package cashcontrol

import (
	"context"
	"fmt"

	cashcontrolmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/cashcontrol"
	cashcontrolmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/app/cashcontrol"
)

func (h *Handler) UpdateCashControl(ctx context.Context) (*cashcontrolmwpb.CashControl, error) {
	info, err := h.GetCashControl(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	if info.AppID != *h.AppID || info.EntID != *h.EntID {
		return nil, fmt.Errorf("permission denied")
	}

	if _, err := cashcontrolmwcli.UpdateCashControl(
		ctx,
		&cashcontrolmwpb.CashControlReq{
			ID:    h.ID,
			Value: h.Value,
		},
	); err != nil {
		return nil, err
	}
	return h.GetCashControl(ctx)
}
