package coin

import (
	"context"
	"fmt"

	eventcoinmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coin"
	eventcoinmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event/coin"
)

func (h *Handler) UpdateEventCoin(ctx context.Context) (*npool.EventCoin, error) {
	info, err := h.GetEventCoin(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid event")
	}
	if info.ID != *h.ID || info.EntID != *h.EntID || info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	if err := eventcoinmwcli.UpdateEventCoin(ctx, &eventcoinmwpb.EventCoinReq{
		ID:         h.ID,
		EntID:      h.EntID,
		AppID:      h.AppID,
		CoinValue:  h.CoinValue,
		CoinPreUSD: h.CoinPreUSD,
	}); err != nil {
		return nil, err
	}
	return h.GetEventCoin(ctx)
}
