package config

import (
	"context"

	"github.com/google/uuid"

	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
)

func (h *Handler) CreateAppConfig(ctx context.Context) (*npool.AppConfig, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	if _, err := appconfigmwcli.CreateAppConfig(ctx, &appconfigmwpb.AppConfigReq{
		EntID:            h.EntID,
		AppID:            h.AppID,
		CommissionType:   h.CommissionType,
		SettleMode:       h.SettleMode,
		SettleAmountType: h.SettleAmountType,
		SettleInterval:   h.SettleInterval,
		StartAt:          h.StartAt,
		SettleBenefit:    h.SettleBenefit,
		MaxLevel:         h.MaxLevel,
	}); err != nil {
		return nil, err
	}

	return h.GetAppConfig(ctx)
}
