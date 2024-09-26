package coin

import (
	"context"
	"fmt"

	coinconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	eventcoinmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event/coin"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coin"
	coinconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	eventcoinmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event/coin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) checkCoinConfig(ctx context.Context) error {
	info, err := coinconfigmwcli.GetCoinConfigOnly(ctx, &coinconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinConfigID},
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid coin")
	}

	maxValue, err := decimal.NewFromString(info.MaxValue)
	if err != nil {
		return err
	}
	coinValue, err := decimal.NewFromString(*h.CoinValue)
	if err != nil {
		return err
	}
	if coinValue.Cmp(maxValue) > 0 {
		return fmt.Errorf("invalid coinvalue")
	}
	if h.CoinPerUSD != nil {
		coinPerUSD, err := decimal.NewFromString(*h.CoinPerUSD)
		if err != nil {
			return err
		}
		if coinPerUSD.Cmp(maxValue) > 0 {
			return fmt.Errorf("invalid coinperusd")
		}
	}

	return nil
}

func (h *createHandler) checkEvent(ctx context.Context) error {
	exist, err := eventmwcli.ExistEventConds(ctx, &eventmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EventID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid event")
	}

	return nil
}

func (h *Handler) CreateEvent(ctx context.Context) (*npool.EventCoin, error) {
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.checkEvent(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkCoinConfig(ctx); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	req := &eventcoinmwpb.EventCoinReq{
		EntID:        h.EntID,
		AppID:        h.AppID,
		EventID:      h.EventID,
		CoinConfigID: h.CoinConfigID,
		CoinValue:    h.CoinValue,
		CoinPerUSD:   h.CoinPerUSD,
	}

	if err := eventcoinmwcli.CreateEventCoin(ctx, req); err != nil {
		return nil, err
	}

	return h.GetEventCoin(ctx)
}
