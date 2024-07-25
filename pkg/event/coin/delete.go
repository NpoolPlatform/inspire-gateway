package coin

import (
	"context"
	"fmt"

	eventcoinmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event/coin"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coin"
	eventcoinmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event/coin"
)

func (h *Handler) DeleteEventCoin(ctx context.Context) (*npool.EventCoin, error) {
	info, err := eventcoinmwcli.GetEventCoinOnly(ctx, &eventcoinmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid eventcoin")
	}

	if err := eventcoinmwcli.DeleteEventCoin(ctx, h.ID, h.EntID); err != nil {
		return nil, err
	}

	return h.GetEventCoinExt(ctx, info)
}
