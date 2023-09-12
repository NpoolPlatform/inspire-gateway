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

func (h *updateHandler) checkAppGood(ctx context.Context) error {
	if h.AppGoodID == nil {
		return nil
	}
	if h.AppID == nil {
		return fmt.Errorf("invalid appid")
	}

	good, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
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
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	handler := &updateHandler{
		Handler: h,
	}

	if err := handler.checkAppGood(ctx); err != nil {
		return nil, err
	}
	_, err := eventmwcli.UpdateEvent(ctx, &eventmwpb.EventReq{
		ID:             h.ID,
		AppID:          h.AppID,
		CouponIDs:      h.CouponIDs,
		Credits:        h.Credits,
		CreditsPerUSD:  h.CreditsPerUSD,
		MaxConsecutive: h.MaxConsecutive,
		GoodID:         &handler.appGood.GoodID,
		AppGoodID:      h.AppGoodID,
		InviterLayers:  h.InviterLayers,
	})
	if err != nil {
		return nil, err
	}
	return h.GetEvent(ctx)
}
