package config

import (
	"context"
	"fmt"

	coinconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	coinconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
)

func (h *Handler) GetCoinConfig(ctx context.Context) (*coinconfigmwpb.CoinConfig, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}

	info, err := coinconfigmwcli.GetCoinConfig(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	return info, nil
}

func (h *Handler) GetCoinConfigs(ctx context.Context) ([]*coinconfigmwpb.CoinConfig, uint32, error) {
	conds := &coinconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}

	return coinconfigmwcli.GetCoinConfigs(ctx, conds, h.Offset, h.Limit)
}
