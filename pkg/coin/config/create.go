package config

import (
	"context"

	"github.com/google/uuid"

	coinconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/config"
	coinconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
)

func (h *Handler) CreateCoinConfig(ctx context.Context) (*npool.CoinConfig, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	if err := coinconfigmwcli.CreateCoinConfig(ctx, &coinconfigmwpb.CoinConfigReq{
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
