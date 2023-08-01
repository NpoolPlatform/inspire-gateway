package event

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/appgood"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
)

type queryHandler struct {
	*Handler
	events  []*eventmwpb.Event
	app     *appmwpb.App
	goods   map[string]*appgoodmwpb.Good
	coupons map[string]*couponmwpb.Coupon
	infos   []*npool.Event
}

func (h *queryHandler) getApp(ctx context.Context) error {
	if h.AppID == nil {
		return fmt.Errorf("invalid appid")
	}
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

func (h *queryHandler) getGoods(ctx context.Context) error {
	goodIDs := []string{}
	for _, event := range h.events {
		if event.GoodID != nil {
			goodIDs = append(goodIDs, *event.GoodID)
		}
	}
	goods, _, err := appgoodmwcli.GetGoods(
		ctx,
		&appgoodmgrpb.Conds{
			AppID:   &commonpb.StringVal{Op: cruder.EQ, Value: *h.AppID},
			GoodIDs: &commonpb.StringSliceVal{Op: cruder.IN, Value: goodIDs},
		},
		0,
		int32(len(goodIDs)),
	)
	if err != nil {
		return err
	}
	for _, good := range goods {
		h.goods[good.GoodID] = good
	}
	return nil
}

func (h *queryHandler) getCoupons(ctx context.Context) error {
	couponIDs := []string{}
	for _, event := range h.events {
		couponIDs = append(couponIDs, event.CouponIDs...)
	}
	coupons, _, err := couponmwcli.GetCoupons(
		ctx,
		&couponmwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			IDs:   &basetypes.StringSliceVal{Op: cruder.IN, Value: couponIDs},
		},
		0,
		int32(len(couponIDs)),
	)
	if err != nil {
		return err
	}
	for _, coupon := range coupons {
		h.coupons[coupon.ID] = coupon
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, event := range h.events {
		info := &npool.Event{
			ID:             event.ID,
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

		if event.GoodID != nil {
			if good, ok := h.goods[*event.GoodID]; ok {
				info.GoodID = *event.GoodID
				info.GoodName = good.GoodName
			}
		}

		for _, couponID := range event.CouponIDs {
			coupon, ok := h.coupons[couponID]
			if !ok {
				continue
			}
			info.Coupons = append(info.Coupons, coupon)
		}

		h.infos = append(h.infos, info)
	}
}

func (h *Handler) GetEvent(ctx context.Context) (*npool.Event, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := eventmwcli.GetEvent(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler: h,
		events:  []*eventmwpb.Event{info},
	}
	handler.AppID = &info.AppID
	if err := handler.getApp(ctx); err != nil {
		return nil, err
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, err
	}
	if err := handler.getCoupons(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetEvents(ctx context.Context) ([]*npool.Event, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appid")
	}

	infos, total, err := eventmwcli.GetEvents(
		ctx,
		&eventmwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		},
		h.Offset,
		h.Limit,
	)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	handler := &queryHandler{
		Handler: h,
		events:  infos,
	}
	if err := handler.getApp(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCoupons(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, total, nil
	}
	return handler.infos, total, nil
}
