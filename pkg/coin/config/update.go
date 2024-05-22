package config

import (
	"context"
	"fmt"

	configmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/config"
	configmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
)

func (h *Handler) UpdateCoinConfig(ctx context.Context) (*npool.CoinConfig, error) {
	info, err := h.GetCoinConfig(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid config")
	}
	if info.ID != *h.ID || info.EntID != *h.EntID {
		return nil, fmt.Errorf("permission denied")
	}

	if err := configmwcli.UpdateCoinConfig(ctx, &configmwpb.CoinConfigReq{
		ID:         h.ID,
		EntID:      h.EntID,
		AppID:      h.AppID,
		CoinTypeID: h.CoinTypeID,
		MaxValue:   h.MaxValue,
		Allocated:  h.Allocated,
	}); err != nil {
		return nil, err
	}
	return h.GetCoinConfig(ctx)
}
