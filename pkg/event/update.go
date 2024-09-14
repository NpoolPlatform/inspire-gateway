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

type updateHandler struct {
	*Handler
	appGood *appgoodmwpb.Good
}

//nolint:dupl
func (h *updateHandler) checkAppGood(ctx context.Context) error {
	if h.AppGoodID == nil {
		return nil
	}

	good, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
	})
	if err != nil {
		return err
	}
	if good == nil {
		return fmt.Errorf("invalid goodid")
	}

	h.appGood = good
	return nil
}

func (h *Handler) UpdateEvent(ctx context.Context) (*npool.Event, error) {
	info, err := h.GetEvent(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid event")
	}
	if info.ID != *h.ID || info.EntID != *h.EntID || info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	handler := &updateHandler{
		Handler: h,
	}
	if err := handler.checkAppGood(ctx); err != nil {
		return nil, err
	}

	req := &eventmwpb.EventReq{
		ID:             h.ID,
		EntID:          h.EntID,
		AppID:          h.AppID,
		Credits:        h.Credits,
		CreditsPerUSD:  h.CreditsPerUSD,
		MaxConsecutive: h.MaxConsecutive,
		InviterLayers:  h.InviterLayers,
	}
	if handler.appGood != nil {
		req.GoodID = &handler.appGood.GoodID
		req.AppGoodID = h.AppGoodID
	}

	if _, err := eventmwcli.UpdateEvent(ctx, req); err != nil {
		return nil, err
	}
	return h.GetEvent(ctx)
}
