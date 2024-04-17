package config

import (
	"context"

	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/commission/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/commission/config"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/commission/config"
)

type updateHandler struct {
	*checkHandler
}

func (h *Handler) UpdateCommission(ctx context.Context) (*npool.AppCommissionConfig, error) {
	handler := &updateHandler{
		checkHandler: &checkHandler{
			Handler: h,
		},
	}
	if err := handler.checkConfig(ctx); err != nil {
		return nil, err
	}

	_, err := commissionconfigmwcli.UpdateCommissionConfig(ctx, &commissionconfigmwpb.AppCommissionConfigReq{
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
