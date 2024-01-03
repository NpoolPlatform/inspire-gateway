package cashcontrol

import (
	"context"
	"fmt"

	cashcontrolmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/cashcontrol"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/app/cashcontrol"
)

func (h *Handler) GetCashControl(ctx context.Context) (*npool.CashControl, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}
	info, err := cashcontrolmwcli.GetCashControl(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	return info, nil
}

func (h *Handler) GetCashControls(ctx context.Context) ([]*npool.CashControl, uint32, error) {
	conds := &npool.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	return cashcontrolmwcli.GetCashControls(ctx, conds, h.Offset, h.Limit)
}
