package coupon

import (
	"context"
	"fmt"

	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	eventcouponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event/coupon"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coupon"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	eventcouponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event/coupon"
)

type queryHandler struct {
	*Handler
	eventCoupons []*eventcouponmwpb.EventCoupon
	coupons      map[string]*couponmwpb.Coupon
	infos        []*npool.EventCoupon
	events       map[string]*eventmwpb.Event
}

func (h *queryHandler) getEvents(ctx context.Context) error {
	infos, _, err := eventmwcli.GetEvents(ctx, &eventmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return err
	}
	for _, event := range infos {
		h.events[event.EntID] = event
	}
	return nil
}

func (h *queryHandler) getAppCoins(ctx context.Context) error {
	couponIDs := []string{}
	for _, val := range h.eventCoupons {
		couponIDs = append(couponIDs, val.CouponID)
	}
	coupons, _, err := couponmwcli.GetCoupons(
		ctx,
		&couponmwpb.Conds{
			AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: couponIDs},
		},
		0,
		int32(len(couponIDs)),
	)
	if err != nil {
		return err
	}

	for _, coupon := range coupons {
		h.coupons[coupon.EntID] = coupon
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, eventCoupon := range h.eventCoupons {
		coupon, ok := h.coupons[eventCoupon.CouponID]
		if !ok {
			continue
		}
		_, ok = h.events[eventCoupon.EventID]
		if !ok {
			continue
		}

		h.infos = append(h.infos, &npool.EventCoupon{
			ID:           eventCoupon.ID,
			EntID:        eventCoupon.EntID,
			AppID:        eventCoupon.AppID,
			EventID:      eventCoupon.EventID,
			CouponID:     eventCoupon.CouponID,
			CouponType:   coupon.CouponType,
			Denomination: coupon.Denomination,
			Circulation:  coupon.Circulation,
			StartAt:      coupon.StartAt,
			EndAt:        coupon.EndAt,
			DurationDays: coupon.DurationDays,
			Name:         coupon.Name,
		})
	}
}

func (h *Handler) GetEventCoupon(ctx context.Context) (*npool.EventCoupon, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := eventcouponmwcli.GetEventCoupon(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:      h,
		eventCoupons: []*eventcouponmwpb.EventCoupon{info},
		coupons:      map[string]*couponmwpb.Coupon{},
		events:       map[string]*eventmwpb.Event{},
	}
	handler.AppID = &info.AppID
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, err
	}
	if err := handler.getEvents(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetEventCouponExt(ctx context.Context, info *eventcouponmwpb.EventCoupon) (*npool.EventCoupon, error) {
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:      h,
		eventCoupons: []*eventcouponmwpb.EventCoupon{info},
		coupons:      map[string]*couponmwpb.Coupon{},
		events:       map[string]*eventmwpb.Event{},
	}
	handler.AppID = &info.AppID
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, err
	}
	if err := handler.getEvents(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetEventCoupons(ctx context.Context) ([]*npool.EventCoupon, uint32, error) {
	conds := &eventcouponmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.EventID != nil {
		conds.EventID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.EventID}
	}
	infos, total, err := eventcouponmwcli.GetEventCoupons(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	handler := &queryHandler{
		Handler:      h,
		eventCoupons: infos,
		coupons:      map[string]*couponmwpb.Coupon{},
		events:       map[string]*eventmwpb.Event{},
	}
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getEvents(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, total, nil
	}
	return handler.infos, total, nil
}
