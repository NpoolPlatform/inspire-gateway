package event

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
)

type queryHandler struct {
	*Handler
	events   []*eventmwpb.Event
	app      *appmwpb.App
	appGoods map[string]*appgoodmwpb.Good
	infos    []*npool.Event
}

func (h *queryHandler) getApp(ctx context.Context) error {
	app, err := appmwcli.GetApp(ctx, *h.AppID)
	if err != nil {
		return err
	}
	if app == nil {
		return fmt.Errorf("invalid app")
	}
	h.app = app
	return nil
}

func (h *queryHandler) getAppGoods(ctx context.Context) error {
	goodIDs := []string{}
	for _, event := range h.events {
		if event.AppGoodID != nil {
			goodIDs = append(goodIDs, *event.AppGoodID)
		}
	}
	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: goodIDs},
	}, 0, int32(len(goodIDs)))
	if err != nil {
		return err
	}
	for _, good := range goods {
		h.appGoods[good.EntID] = good
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, event := range h.events {
		info := &npool.Event{
			ID:             event.ID,
			EntID:          event.EntID,
			AppID:          event.AppID,
			AppName:        h.app.Name,
			EventType:      event.EventType,
			Credits:        event.Credits,
			CreditsPerUSD:  event.CreditsPerUSD,
			MaxConsecutive: event.MaxConsecutive,
			InviterLayers:  event.InviterLayers,
			CreatedAt:      event.CreatedAt,
			UpdatedAt:      event.UpdatedAt,
		}

		if event.AppGoodID != nil {
			if good, ok := h.appGoods[*event.AppGoodID]; ok {
				info.GoodID = *event.GoodID
				info.AppGoodID = *event.AppGoodID
				info.GoodName = good.GoodName
				info.AppGoodName = good.AppGoodName
			}
		}

		h.infos = append(h.infos, info)
	}
}

func (h *Handler) GetEvent(ctx context.Context) (*npool.Event, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := eventmwcli.GetEvent(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:  h,
		events:   []*eventmwpb.Event{info},
		appGoods: map[string]*appgoodmwpb.Good{},
	}
	if h.AppID == nil {
		handler.AppID = &info.AppID
	}
	if err := handler.getApp(ctx); err != nil {
		return nil, err
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetEventExt(ctx context.Context, info *eventmwpb.Event) (*npool.Event, error) {
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:  h,
		events:   []*eventmwpb.Event{info},
		appGoods: map[string]*appgoodmwpb.Good{},
	}
	handler.AppID = &info.AppID
	if err := handler.getApp(ctx); err != nil {
		return nil, err
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetEvents(ctx context.Context) ([]*npool.Event, uint32, error) {
	infos, total, err := eventmwcli.GetEvents(ctx, &eventmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	handler := &queryHandler{
		Handler:  h,
		events:   infos,
		appGoods: map[string]*appgoodmwpb.Good{},
	}
	if err := handler.getApp(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, total, nil
	}
	return handler.infos, total, nil
}
