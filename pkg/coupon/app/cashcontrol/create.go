package cashcontrol

import (
	"context"
	"fmt"

	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	cashcontrolmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/cashcontrol"
	cashcontrolmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/app/cashcontrol"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) getCoupon(ctx context.Context) error {
	coupon, err := couponmwcli.GetCoupon(ctx, *h.CouponID)
	if err != nil {
		return err
	}
	if coupon == nil {
		return fmt.Errorf("coupon not exist")
	}
	h.AppID = &coupon.AppID
	return nil
}

func (h *createHandler) createCashControl(ctx context.Context) error {
	if _, err := cashcontrolmwcli.CreateCashControl(
		ctx,
		&cashcontrolmwpb.CashControlReq{
			EntID:       h.EntID,
			AppID:       h.AppID,
			CouponID:    h.CouponID,
			ControlType: h.ControlType,
			Value:       h.Value,
		},
	); err != nil {
		return err
	}
	return nil
}

func (h *Handler) CreateCashControl(ctx context.Context) (*cashcontrolmwpb.CashControl, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	handler := &createHandler{
		Handler: h,
	}
	if err := handler.getCoupon(ctx); err != nil {
		return nil, err
	}
	if err := handler.createCashControl(ctx); err != nil {
		return nil, err
	}

	return h.GetCashControl(ctx)
}
