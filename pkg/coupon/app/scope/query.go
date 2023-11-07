package scope

import (
	"context"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	appgoodscopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/scope"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/scope"
	appgoodscopemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/app/scope"
)

type queryHandler struct {
	*Handler
	infos         []*npool.Scope
	appgoodscopes []*appgoodscopemwpb.Scope
	goods         map[string]*goodmwpb.Good
	appgoods      map[string]*appgoodmwpb.Good
}

func (h *queryHandler) getGoods(ctx context.Context) error {
	ids := []string{}
	for _, info := range h.appgoodscopes {
		ids = append(ids, info.GoodID)
	}

	goods, _, err := goodmwcli.GetGoods(ctx, &goodmwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, int32(0), int32(len(ids)))
	if err != nil {
		return err
	}
	for _, good := range goods {
		h.goods[good.ID] = good
	}
	return nil
}

func (h *queryHandler) getAppGoods(ctx context.Context) error {
	ids := []string{}
	for _, info := range h.appgoodscopes {
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
	for _, info := range h.appgoodscopes {
		_, ok := h.goods[info.GoodID]
		if !ok {
			continue
		}
		appgood, ok := h.appgoods[info.AppGoodID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.Scope{
			ID:                 info.ID,
			AppID:              info.AppID,
			AppGoodID:          info.AppGoodID,
			GoodName:           appgood.GoodName,
			ScopeID:            info.ScopeID,
			GoodID:             info.GoodID,
			CouponID:           info.CouponID,
			CouponName:         info.CouponName,
			CouponType:         info.CouponType,
			CouponScope:        info.CouponScope,
			CouponDenomination: info.CouponDenomination,
			CreatedAt:          info.CreatedAt,
			UpdatedAt:          info.UpdatedAt,
		})
	}
}

func (h *Handler) GetAppGoodScope(ctx context.Context) (*npool.Scope, error) {
	info, err := appgoodscopemwcli.GetAppGoodScope(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:       h,
		appgoodscopes: []*appgoodscopemwpb.Scope{info},
		goods:         map[string]*goodmwpb.Good{},
		appgoods:      map[string]*appgoodmwpb.Good{},
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, err
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

func (h *Handler) GetAppGoodScopes(ctx context.Context) ([]*npool.Scope, uint32, error) {
	scopes, total, err := appgoodscopemwcli.GetAppGoodScopes(ctx, &appgoodscopemwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}

	handler := &queryHandler{
		Handler:       h,
		appgoodscopes: scopes,
		goods:         map[string]*goodmwpb.Good{},
		appgoods:      map[string]*appgoodmwpb.Good{},
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, 0, err
	}

	handler.formalize()
	return handler.infos, total, nil
}
