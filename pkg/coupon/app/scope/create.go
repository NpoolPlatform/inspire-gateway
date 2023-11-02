package scope

import (
	"context"
	"fmt"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	scopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/scope"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"
	scopemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/scope"

	appgoodscopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/scope"
	appgoodscopemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/app/scope"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
	scope   *scopemwpb.Scope
	good    *goodmwpb.Good
	appgood *appgoodmwpb.Good
}

func (h *createHandler) getScope(ctx context.Context) error {
	scope, err := scopemwcli.GetScope(ctx, *h.ScopeID)
	if err != nil {
		return err
	}
	if scope == nil {
		return fmt.Errorf("scopeid not exist")
	}
	h.scope = scope
	return nil
}

func (h *createHandler) getGood(ctx context.Context) error {
	good, err := goodmwcli.GetGood(ctx, h.scope.GoodID)
	if err != nil {
		return err
	}
	if good == nil {
		return fmt.Errorf("good not exist")
	}
	h.good = good
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
	if h.scope.GoodID != appgood.GoodID {
		return fmt.Errorf("goodid mismatch")
	}
	h.appgood = appgood
	return nil
}
func (h *createHandler) createAppGoodScope(ctx context.Context) error {
	if h.CouponScope == nil {
		h.CouponScope = &h.scope.CouponScope
	}

	if _, err := appgoodscopemwcli.CreateAppGoodScope(
		ctx,
		&appgoodscopemwpb.ScopeReq{
			ID:          h.ID,
			AppID:       h.AppID,
			AppGoodID:   h.AppGoodID,
			ScopeID:     h.ScopeID,
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
	if err := handler.getScope(ctx); err != nil {
		return nil, err
	}
	if err := handler.getGood(ctx); err != nil {
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
