//nolint:dupl
package event

import (
	"context"
	"fmt"

	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
)

func (h *Handler) UpdateEvent(ctx context.Context) (*npool.Event, error) {
	info, err := h.GetEvent(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid event")
	}
	if info.ID != *h.ID || info.EntID != *h.EntID {
		return nil, fmt.Errorf("permission denied")
	}

	if _, err := eventmwcli.UpdateEvent(ctx, &eventmwpb.EventReq{
		ID:             h.ID,
		EntID:          h.EntID,
		AppID:          h.AppID,
		CouponIDs:      h.CouponIDs,
		Credits:        h.Credits,
		CreditsPerUSD:  h.CreditsPerUSD,
		MaxConsecutive: h.MaxConsecutive,
		InviterLayers:  h.InviterLayers,
	}); err != nil {
		return nil, err
	}
	return h.GetEvent(ctx)
}
