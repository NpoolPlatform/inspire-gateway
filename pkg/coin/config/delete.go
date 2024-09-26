package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
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
		return nil, wlog.WrapError(err)
	}
	if info == nil {
		return nil, wlog.Errorf("invalid coinconfig")
	}
	h.AppID = &info.AppID

	if err := configmwcli.DeleteCoinConfig(ctx, h.ID, h.EntID); err != nil {
		return nil, wlog.WrapError(err)
	}

	return h.GetCoinConfigExt(ctx, info)
}
