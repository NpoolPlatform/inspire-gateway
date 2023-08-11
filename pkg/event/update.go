package event

import (
	"context"
	"fmt"

	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
)

func (h *Handler) UpdateEvent(ctx context.Context) (*npool.Event, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	_, err := eventmwcli.UpdateEvent(ctx, &eventmwpb.EventReq{
		ID:             h.ID,
		AppID:          h.AppID,
		CouponIDs:      h.CouponIDs,
		Credits:        h.Credits,
		CreditsPerUSD:  h.CreditsPerUSD,
		MaxConsecutive: h.MaxConsecutive,
		GoodID:         h.GoodID,
		InviterLayers:  h.InviterLayers,
	})
	if err != nil {
		return nil, err
	}
	return h.GetEvent(ctx)
}
