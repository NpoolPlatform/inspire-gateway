package config

import (
	"context"
	"fmt"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	coinconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/config"
	coinconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
)

type queryHandler struct {
	*Handler
	appcoin     map[string]*appcoinmwpb.Coin
	coinConfigs []*coinconfigmwpb.CoinConfig
	infos       []*npool.CoinConfig
}

func (h *queryHandler) getAppCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, val := range h.coinConfigs {
		coinTypeIDs = append(coinTypeIDs, val.CoinTypeID)
	}
	coins, _, err := appcoinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		AppID:       &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: coinTypeIDs},
	}, 0, int32(len(coinTypeIDs)))
	if err != nil {
		return err
	}

	for _, coin := range coins {
		h.appcoin[coin.CoinTypeID] = coin
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, coinConfig := range h.coinConfigs {
		appcoin, ok := h.appcoin[coinConfig.CoinTypeID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.CoinConfig{
			ID:           coinConfig.ID,
			EntID:        coinConfig.EntID,
			AppID:        coinConfig.AppID,
			CoinTypeID:   coinConfig.CoinTypeID,
			MaxValue:     coinConfig.MaxValue,
			Allocated:    coinConfig.Allocated,
			CoinName:     appcoin.CoinName,
			DisplayNames: appcoin.DisplayNames,
			CoinLogo:     appcoin.Logo,
			CoinUnit:     appcoin.Unit,
			CreatedAt:    coinConfig.CreatedAt,
			UpdatedAt:    coinConfig.UpdatedAt,
		})
	}
}

func (h *Handler) GetCoinConfig(ctx context.Context) (*npool.CoinConfig, error) {
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
	handler := &queryHandler{
		Handler:     h,
		appcoin:     map[string]*appcoinmwpb.Coin{},
		coinConfigs: []*coinconfigmwpb.CoinConfig{info},
	}
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}

	return handler.infos[0], nil
}

func (h *Handler) GetCoinConfigs(ctx context.Context) ([]*npool.CoinConfig, uint32, error) {
	conds := &coinconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}

	infos, _, err := coinconfigmwcli.GetCoinConfigs(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	handler := &queryHandler{
		Handler:     h,
		appcoin:     map[string]*appcoinmwpb.Coin{},
		coinConfigs: []*coinconfigmwpb.CoinConfig{},
	}

	handler.coinConfigs = infos
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()

	return handler.infos, uint32(len(handler.infos)), nil
}
