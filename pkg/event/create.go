package event

import (
	"context"

	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
)

func (h *Handler) CreateEvent(ctx context.Context) (*npool.Event, error) {
	info, err := eventmwcli.CreateEvent(ctx, &eventmwpb.EventReq{
		AppID:          h.AppID,
		EventType:      h.EventType,
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
	h.ID = &info.ID
	return h.GetEvent(ctx)
}
