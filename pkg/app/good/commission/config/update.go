package config

import (
	"context"
	"fmt"

	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/good/commission/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/good/commission/config"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/good/commission/config"
)

func (h *Handler) UpdateCommission(ctx context.Context) (*npool.AppGoodCommissionConfig, error) {
	info, err := commissionconfigmwcli.GetCommissionConfig(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid commission")
	}
	if info.ID != *h.ID || info.AppID != *h.AppID || info.EndAt != 0 {
		return nil, fmt.Errorf("permission denied")
	}

	_, err = commissionconfigmwcli.UpdateCommissionConfig(ctx, &commissionconfigmwpb.AppGoodCommissionConfigReq{
		ID:              h.ID,
		StartAt:         h.StartAt,
		ThresholdAmount: h.ThresholdAmount,
		Invites:         h.Invites,
		Disabled:        h.Disabled,
		Level:           h.Level,
	})
	if err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
