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
)

type createHandler struct {
	*Handler
}

func (h *createHandler) checkCoinConfig(ctx context.Context) error {
	exist, err := coinconfigmwcli.ExistCoinConfigConds(ctx, &coinconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinConfigID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid coin")
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
		CoinPreUSD:   h.CoinPreUSD,
	}

	if err := eventcoinmwcli.CreateEventCoin(ctx, req); err != nil {
		return nil, err
	}

	return h.GetEventCoin(ctx)
}
