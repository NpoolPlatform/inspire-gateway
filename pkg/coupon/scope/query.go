package scope

import (
	"context"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	scopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/scope"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"
	scopemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/scope"
)

type queryHandler struct {
	*Handler
	infos    []*npool.Scope
	scopes   []*scopemwpb.Scope
	appgoods map[string]*appgoodmwpb.Good
}

func (h *queryHandler) getAppGoods(ctx context.Context) error {
	ids := []string{}
	for _, info := range h.scopes {
		ids = append(ids, info.AppGoodID)
	}
	appgoods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		IDs:   &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, int32(0), int32(len(ids)))
	if err != nil {
		return err
	}
	for _, appgood := range appgoods {
		h.appgoods[appgood.ID] = appgood
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, info := range h.scopes {
		appgood, ok := h.appgoods[info.AppGoodID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.Scope{
			ID:                 info.ID,
			AppGoodID:          info.AppGoodID,
			GoodName:           appgood.GoodName,
			CouponID:           info.CouponID,
			CouponName:         info.CouponName,
			CouponScope:        info.CouponScope,
			CouponScopeStr:     info.CouponScopeStr,
			CouponType:         info.CouponType,
			CouponTypeStr:      info.CouponTypeStr,
			CouponDenomination: info.CouponDenomination,
			CreatedAt:          info.CreatedAt,
			UpdatedAt:          info.UpdatedAt,
		})
	}
}

func (h *Handler) GetScope(ctx context.Context) (*npool.Scope, error) {
	info, err := scopemwcli.GetScope(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:  h,
		scopes:   []*scopemwpb.Scope{info},
		appgoods: map[string]*appgoodmwpb.Good{},
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

func (h *Handler) GetScopes(ctx context.Context) ([]*npool.Scope, uint32, error) {
	conds := &scopemwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	infos, total, err := scopemwcli.GetScopes(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	handler := &queryHandler{
		Handler:  h,
		scopes:   infos,
		appgoods: map[string]*appgoodmwpb.Good{},
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, 0, err
	}

	handler.formalize()
	return handler.infos, total, nil
}
