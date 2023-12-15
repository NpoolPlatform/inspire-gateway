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
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
	appGood *appgoodmwpb.Good
}

func (h *createHandler) checkAppGood(ctx context.Context) error {
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

func (h *Handler) CreateEvent(ctx context.Context) (*npool.Event, error) {
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.checkAppGood(ctx); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	req := &eventmwpb.EventReq{
		EntID:          h.EntID,
		AppID:          h.AppID,
		EventType:      h.EventType,
		CouponIDs:      h.CouponIDs,
		Credits:        h.Credits,
		CreditsPerUSD:  h.CreditsPerUSD,
		MaxConsecutive: h.MaxConsecutive,
		InviterLayers:  h.InviterLayers,
	}
	if handler.appGood != nil {
		req.GoodID = &handler.appGood.GoodID
		req.AppGoodID = h.AppGoodID
	}

	if _, err := eventmwcli.CreateEvent(ctx, req); err != nil {
		return nil, err
	}

	return h.GetEvent(ctx)
}
