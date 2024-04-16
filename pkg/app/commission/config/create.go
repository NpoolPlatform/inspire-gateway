package config

import (
	"context"

	"github.com/google/uuid"

	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/commission/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/commission/config"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/commission/config"
)

func (h *Handler) CreateCommissionConfig(ctx context.Context) (*npool.AppCommissionConfig, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	if _, err := commissionconfigmwcli.CreateCommissionConfig(ctx, &commissionconfigmwpb.AppCommissionConfigReq{
		EntID:           h.EntID,
		AppID:           h.AppID,
		SettleType:      h.SettleType,
		Invites:         h.Invites,
		StartAt:         h.StartAt,
		AmountOrPercent: h.AmountOrPercent,
		ThresholdAmount: h.ThresholdAmount,
		Disabled:        h.Disabled,
		Level:           h.Level,
	}); err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
