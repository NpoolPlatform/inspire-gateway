package scope

import (
	"context"
	"fmt"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	scopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/scope"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"
	scopemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/scope"
)

type queryHandler struct {
	*Handler
	infos  []*npool.Scope
	scopes []*scopemwpb.Scope
	goods  map[string]*goodmwpb.Good
}

func (h *queryHandler) getGoods(ctx context.Context) error {
	ids := []string{}
	for _, info := range h.scopes {
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

func (h *queryHandler) formalize() {
	for _, info := range h.scopes {
		good, ok := h.goods[info.GoodID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.Scope{
			ID:                 info.ID,
			EntID:              info.EntID,
			GoodID:             info.GoodID,
			GoodTitle:          good.Title,
			CouponID:           info.CouponID,
			CouponName:         info.CouponName,
			CouponType:         info.CouponType,
			CouponScope:        info.CouponScope,
			CouponDenomination: info.CouponDenomination,
			CouponCirculation:  info.CouponCirculation,
			CreatedAt:          info.CreatedAt,
			UpdatedAt:          info.UpdatedAt,
		})
	}
}

func (h *Handler) GetScope(ctx context.Context) (*npool.Scope, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}
	info, err := scopemwcli.GetScope(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler: h,
		scopes:  []*scopemwpb.Scope{info},
		goods:   map[string]*goodmwpb.Good{},
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, err
	}

	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetScopes(ctx context.Context) ([]*npool.Scope, uint32, error) {
	scopes, total, err := scopemwcli.GetScopes(ctx, &scopemwpb.Conds{}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}

	handler := &queryHandler{
		Handler: h,
		scopes:  scopes,
		goods:   map[string]*goodmwpb.Good{},
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, 0, err
	}

	handler.formalize()
	return handler.infos, total, nil
}
