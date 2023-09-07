package event

import (
	"context"
	"fmt"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) validateGood(ctx context.Context) error {
	if h.GoodID == nil {
		return nil
	}
	if h.AppID == nil {
		return fmt.Errorf("invalid appid")
	}

	good, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		GoodID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.GoodID},
	})
	if err != nil {
		return err
	}
	if good == nil {
		return fmt.Errorf("invalid goodid")
	}

	return nil
}

func (h *Handler) CreateEvent(ctx context.Context) (*npool.Event, error) {
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.validateGood(ctx); err != nil {
		return nil, err
	}

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
