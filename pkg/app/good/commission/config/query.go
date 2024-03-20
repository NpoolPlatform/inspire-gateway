package config

import (
	"context"
	"fmt"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	commconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/good/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/good/commission/config"
	commconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/good/commission/config"

	"github.com/google/uuid"
)

type queryHandler struct {
	*Handler
	appGoods map[string]*appgoodmwpb.Good
	goods    map[string]*goodmwpb.Good
	coins    map[string]*appcoinmwpb.Coin
	comms    []*commconfigmwpb.AppGoodCommissionConfig
	infos    []*npool.AppGoodCommissionConfig
}

func (h *queryHandler) getAppGoods(ctx context.Context) error {
	goodIDs := []string{}
	for _, comm := range h.comms {
		if _, err := uuid.Parse(comm.AppGoodID); err != nil {
			continue
		}
		goodIDs = append(goodIDs, comm.AppGoodID)
	}
	if len(goodIDs) == 0 {
		return nil
	}

	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: h.comms[0].AppID},
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: goodIDs},
	}, int32(0), int32(len(goodIDs)))
	if err != nil {
		return err
	}

	for _, good := range goods {
		h.appGoods[good.EntID] = good
	}
	return nil
}

func (h *queryHandler) getGoods(ctx context.Context) error {
	ids := []string{}
	for _, comm := range h.comms {
		if _, err := uuid.Parse(comm.GoodID); err != nil {
			continue
		}
		ids = append(ids, comm.GoodID)
	}
	if len(ids) == 0 {
		return nil
	}

	goods, _, err := goodmwcli.GetGoods(ctx, &goodmwpb.Conds{
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, int32(0), int32(len(ids)))
	if err != nil {
		return err
	}

	for _, good := range goods {
		h.goods[good.EntID] = good
	}
	return nil
}

func (h *queryHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, good := range h.appGoods {
		coinTypeIDs = append(coinTypeIDs, good.CoinTypeID)
	}
	if len(coinTypeIDs) == 0 {
		return nil
	}

	coins, _, err := coinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		AppID:       &basetypes.StringVal{Op: cruder.EQ, Value: h.comms[0].AppID},
		CoinTypeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: coinTypeIDs},
	}, int32(0), int32(len(coinTypeIDs)))
	if err != nil {
		return err
	}

	for _, coin := range coins {
		h.coins[coin.CoinTypeID] = coin
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, comm := range h.comms {
		good, ok := h.goods[comm.GoodID]
		if !ok {
			continue
		}
		appgood, ok := h.appGoods[comm.AppGoodID]
		if !ok {
			continue
		}
		coin, ok := h.coins[appgood.CoinTypeID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.AppGoodCommissionConfig{
			ID:              comm.ID,
			EntID:           comm.EntID,
			AppID:           comm.AppID,
			SettleType:      comm.SettleType,
			GoodID:          comm.GoodID,
			GoodTitle:       good.Title,
			AppGoodID:       comm.AppGoodID,
			GoodName:        appgood.GoodName,
			AmountOrPercent: comm.AmountOrPercent,
			ThresholdAmount: comm.ThresholdAmount,
			StartAt:         comm.StartAt,
			EndAt:           comm.EndAt,
			CoinTypeID:      appgood.CoinTypeID,
			CoinName:        coin.Name,
			CoinLogo:        coin.Logo,
			CreatedAt:       comm.CreatedAt,
			UpdatedAt:       comm.UpdatedAt,
		})
	}
}

func (h *Handler) GetCommission(ctx context.Context) (*npool.AppGoodCommissionConfig, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}

	info, err := commconfigmwcli.GetCommissionConfig(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:  h,
		goods:    map[string]*goodmwpb.Good{},
		appGoods: map[string]*appgoodmwpb.Good{},
		coins:    map[string]*appcoinmwpb.Coin{},
		comms:    []*commconfigmwpb.AppGoodCommissionConfig{info},
		infos:    []*npool.AppGoodCommissionConfig{},
	}

	if err := handler.getGoods(ctx); err != nil {
		return nil, err
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, err
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}

	return handler.infos[0], nil
}

func (h *Handler) GetCommissions(ctx context.Context) ([]*npool.AppGoodCommissionConfig, uint32, error) {
	handler := &queryHandler{
		Handler:  h,
		goods:    map[string]*goodmwpb.Good{},
		appGoods: map[string]*appgoodmwpb.Good{},
		coins:    map[string]*appcoinmwpb.Coin{},
		infos:    []*npool.AppGoodCommissionConfig{},
	}

	conds := &commconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.EndAt != nil {
		conds.EndAt = &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.EndAt}
	}
	infos, total, err := commconfigmwcli.GetCommissionConfigs(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}
	handler.comms = infos

	if err := handler.getGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()
	return handler.infos, total, nil
}
