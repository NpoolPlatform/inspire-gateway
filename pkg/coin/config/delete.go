package config

import (
	"context"
	"fmt"

	configmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/config"
	configmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
)

func (h *Handler) DeleteCoinConfig(ctx context.Context) (*npool.CoinConfig, error) {
	info, err := configmwcli.GetCoinConfigOnly(ctx, &configmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid coinconfig")
	}

	if err := configmwcli.DeleteCoinConfig(ctx, h.ID, h.EntID); err != nil {
		return nil, err
	}

	return h.GetCoinConfig(ctx)
}
