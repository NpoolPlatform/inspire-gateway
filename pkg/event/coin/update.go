package coin

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	coinconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	eventcoinmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event/coin"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coin"
	coinconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
	eventcoinmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event/coin"
	"github.com/shopspring/decimal"
)

type updateHandler struct {
	*Handler
}

func (h *updateHandler) checkCoinConfig(ctx context.Context) error {
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
	if h.CoinValue != nil {
		coinValue, err := decimal.NewFromString(*h.CoinValue)
		if err != nil {
			return err
		}
		if coinValue.Cmp(maxValue) > 0 {
			return fmt.Errorf("invalid coinvalue")
		}
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

func (h *Handler) UpdateEventCoin(ctx context.Context) (*npool.EventCoin, error) {
	info, err := eventcoinmwcli.GetEventCoin(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid eventcoin")
	}
	if info.ID != *h.ID || info.EntID != *h.EntID || info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	handler := updateHandler{
		Handler: h,
	}

	handler.CoinConfigID = &info.CoinConfigID
	if err := handler.checkCoinConfig(ctx); err != nil {
		return nil, wlog.WrapError(err)
	}

	if err := eventcoinmwcli.UpdateEventCoin(ctx, &eventcoinmwpb.EventCoinReq{
		ID:         h.ID,
		EntID:      h.EntID,
		AppID:      h.AppID,
		CoinValue:  h.CoinValue,
		CoinPerUSD: h.CoinPerUSD,
	}); err != nil {
		return nil, err
	}
	return h.GetEventCoin(ctx)
}
