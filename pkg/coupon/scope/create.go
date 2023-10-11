package scope

import (
	"context"
	"fmt"

	scopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/scope"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"
	scopemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/scope"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) verifyScope(ctx context.Context) error {
	exist, err := scopemwcli.ExistScopeConds(ctx, &scopemwpb.Conds{
		AppID:       &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		AppGoodID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
		CouponID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.CouponID},
		CouponScope: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.CouponScope)},
	})
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("coupon scope already exist")
	}
	return nil
}

func (h *createHandler) createScope(ctx context.Context) error {
	if _, err := scopemwcli.CreateScope(
		ctx,
		&scopemwpb.ScopeReq{
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

func (h *Handler) CreateScope(ctx context.Context) (*npool.Scope, error) {
	if h.AppGoodID == nil {
		switch *h.CouponScope {
		case types.CouponScope_Blacklist:
			fallthrough //nolint
		case types.CouponScope_Whitelist:
			return nil, fmt.Errorf("appgoodid is must")
		}
	}
	appGoodID := uuid.Nil.String()
	if *h.CouponScope == types.CouponScope_AllGood {
		h.AppGoodID = &appGoodID
	}

	id := uuid.NewString()
	if h.ID == nil {
		h.ID = &id
	}

	handler := &createHandler{
		Handler: h,
	}
	if err := handler.verifyScope(ctx); err != nil {
		return nil, err
	}
	if err := handler.createScope(ctx); err != nil {
		return nil, err
	}
	return h.GetScope(ctx)
}
