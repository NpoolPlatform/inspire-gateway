package coin

import (
	"context"
	"fmt"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	eventcoinmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event/coin"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coin"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	eventcoinmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event/coin"
)

type queryHandler struct {
	*Handler
	eventCoins []*eventcoinmwpb.EventCoin
	appcoin    map[string]*appcoinmwpb.Coin
	infos      []*npool.EventCoin
	events     map[string]*eventmwpb.Event
}

func (h *queryHandler) getEvents(ctx context.Context) error {
	infos, _, err := eventmwcli.GetEvents(ctx, &eventmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return wlog.WrapError(err)
	}
	for _, event := range infos {
		h.events[event.EntID] = event
	}
	return nil
}

func (h *queryHandler) getAppCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, val := range h.eventCoins {
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
	for _, eventCoin := range h.eventCoins {
		appcoin, ok := h.appcoin[eventCoin.CoinTypeID]
		if !ok {
			continue
		}
		_, ok = h.events[eventCoin.EventID]
		if !ok {
			continue
		}

		h.infos = append(h.infos, &npool.EventCoin{
			ID:           eventCoin.ID,
			EntID:        eventCoin.EntID,
			AppID:        eventCoin.AppID,
			EventID:      eventCoin.EventID,
			CoinConfigID: eventCoin.CoinConfigID,
			CoinTypeID:   eventCoin.CoinTypeID,
			CoinValue:    eventCoin.CoinValue,
			CoinPerUSD:   eventCoin.CoinPerUSD,
			CoinName:     appcoin.CoinName,
			DisplayNames: appcoin.DisplayNames,
			CoinLogo:     appcoin.Logo,
			CoinUnit:     appcoin.Unit,
		})
	}
}

func (h *Handler) GetEventCoin(ctx context.Context) (*npool.EventCoin, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := eventcoinmwcli.GetEventCoin(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:    h,
		eventCoins: []*eventcoinmwpb.EventCoin{info},
		appcoin:    map[string]*appcoinmwpb.Coin{},
		events:     map[string]*eventmwpb.Event{},
	}
	if h.AppID == nil {
		handler.AppID = &info.AppID
	}
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, err
	}
	if err := handler.getEvents(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetEventCoinExt(ctx context.Context, info *eventcoinmwpb.EventCoin) (*npool.EventCoin, error) {
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:    h,
		eventCoins: []*eventcoinmwpb.EventCoin{info},
		appcoin:    map[string]*appcoinmwpb.Coin{},
		events:     map[string]*eventmwpb.Event{},
	}
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, err
	}
	if err := handler.getEvents(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetEventCoins(ctx context.Context) ([]*npool.EventCoin, uint32, error) {
	conds := &eventcoinmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.EventID != nil {
		conds.EventID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.EventID}
	}
	infos, total, err := eventcoinmwcli.GetEventCoins(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	handler := &queryHandler{
		Handler:    h,
		eventCoins: infos,
		appcoin:    map[string]*appcoinmwpb.Coin{},
		events:     map[string]*eventmwpb.Event{},
	}
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getEvents(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, total, nil
	}
	return handler.infos, total, nil
}
