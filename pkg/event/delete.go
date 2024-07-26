package event

import (
	"context"
	"fmt"

	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
)

func (h *Handler) DeleteEvent(ctx context.Context) (*npool.Event, error) {
	info, err := eventmwcli.GetEventOnly(ctx, &eventmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid taskconfig")
	}

	if err := eventmwcli.DeleteEvent(ctx, h.ID, h.EntID); err != nil {
		return nil, err
	}

	return h.GetEventExt(ctx, info)
}
