package coupon

import (
	"context"
	"fmt"

	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) createSpecialOffer(ctx context.Context) (*couponmwpb.Coupon, error) {
	// TODO: need dtm to create coupon and allocated
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}

func (h *Handler) CreateCoupon(ctx context.Context) (*couponmwpb.Coupon, error) {
	if h.CouponType == nil {
		return nil, fmt.Errorf("invalid coupontype")
	}

	handler := &createHandler{
		Handler: h,
	}

	if *h.CouponType == types.CouponType_SpecialOffer {
		return handler.createSpecialOffer(ctx)
	}
	return couponmwcli.CreateCoupon(
		ctx,
		&couponmwpb.CouponReq{
			ID:               h.ID,
			AppID:            h.AppID,
			CouponType:       h.CouponType,
			Denomination:     h.Denomination,
			Circulation:      h.Circulation,
			IssuedBy:         h.IssuedBy,
			StartAt:          h.StartAt,
			DurationDays:     h.DurationDays,
			Message:          h.Message,
			Name:             h.Name,
			GoodID:           h.GoodID,
			CouponConstraint: h.CouponConstraint,
			Random:           h.Random,
		},
	)
}
