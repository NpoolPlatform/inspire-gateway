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
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	eventcouponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event/coupon"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) checkCoupon(ctx context.Context) error {
	exist, err := couponmwcli.ExistCoupon(ctx, *h.CouponID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid coupon")
	}

	return nil
}

func (h *createHandler) checkEvent(ctx context.Context) error {
	exist, err := eventmwcli.ExistEventConds(ctx, &eventmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EventID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid event")
	}

	return nil
}

func (h *Handler) CreateEvent(ctx context.Context) (*npool.EventCoupon, error) {
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.checkEvent(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkCoupon(ctx); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	req := &eventcouponmwpb.EventCouponReq{
		EntID:    h.EntID,
		AppID:    h.AppID,
		EventID:  h.EventID,
		CouponID: h.CouponID,
	}

	if err := eventcouponmwcli.CreateEventCoupon(ctx, req); err != nil {
		return nil, err
	}

	return h.GetEventCoupon(ctx)
}
