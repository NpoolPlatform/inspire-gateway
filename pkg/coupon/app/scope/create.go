package scope

import (
	"context"
	"fmt"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"

	appgoodscopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/scope"
	appgoodscopemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/app/scope"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
	coupon  *couponmwpb.Coupon
	appgood *appgoodmwpb.Good
}

func (h *createHandler) getCoupon(ctx context.Context) error {
	coupon, err := couponmwcli.GetCoupon(ctx, *h.CouponID)
	if err != nil {
		return err
	}
	if coupon == nil {
		return fmt.Errorf("coupon not exist")
	}
	if coupon.AppID != *h.AppID {
		return fmt.Errorf("permission denied")
	}
	h.coupon = coupon
	return nil
}

func (h *createHandler) getAppGood(ctx context.Context) error {
	appgood, err := appgoodmwcli.GetGood(ctx, *h.AppGoodID)
	if err != nil {
		return err
	}
	if appgood == nil {
		return fmt.Errorf("appgood not exist")
	}
	if appgood.AppID != *h.AppID {
		return fmt.Errorf("permission denied")
	}
	h.appgood = appgood
	return nil
}

func (h *createHandler) createAppGoodScope(ctx context.Context) error {
	if h.CouponScope == nil {
		h.CouponScope = &h.coupon.CouponScope
	}

	if _, err := appgoodscopemwcli.CreateAppGoodScope(
		ctx,
		&appgoodscopemwpb.ScopeReq{
			ID:          h.ID,
			AppID:       h.AppID,
			AppGoodID:   h.AppGoodID,
			CouponID:    h.CouponID,
			CouponScope: h.CouponScope,
		},
	); err != nil {
		return err
	}
	return nil
}

func (h *Handler) CreateAppGoodScope(ctx context.Context) (*npool.Scope, error) {
	id := uuid.NewString()
	if h.ID == nil {
		h.ID = &id
	}

	handler := &createHandler{
		Handler: h,
	}
	if err := handler.getCoupon(ctx); err != nil {
		return nil, err
	}
	if err := handler.getAppGood(ctx); err != nil {
		return nil, err
	}
	if err := handler.createAppGoodScope(ctx); err != nil {
		return nil, err
	}

	return h.GetAppGoodScope(ctx)
}
