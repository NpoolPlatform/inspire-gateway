package coupon

import (
	"context"
	"fmt"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"
)

type updateHandler struct {
	*Handler
	goodID *string
}

func (h *updateHandler) checkGood(ctx context.Context) error {
	if h.AppGoodID == nil {
		return nil
	}
	info, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid good")
	}
	h.goodID = &info.GoodID
	return nil
}

func (h *Handler) UpdateCoupon(ctx context.Context) (*couponmwpb.Coupon, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	handler := &updateHandler{
		Handler: h,
	}
	if err := handler.checkGood(ctx); err != nil {
		return nil, err
	}

	return couponmwcli.UpdateCoupon(ctx, &couponmwpb.CouponReq{
		ID:               h.ID,
		Denomination:     h.Denomination,
		Circulation:      h.Circulation,
		IssuedBy:         h.IssuedBy,
		StartAt:          h.StartAt,
		DurationDays:     h.DurationDays,
		Message:          h.Message,
		Name:             h.Name,
		GoodID:           handler.goodID,
		AppGoodID:        h.AppGoodID,
		CouponConstraint: h.CouponConstraint,
		Threshold:        h.Threshold,
		Random:           h.Random,
	})
}
