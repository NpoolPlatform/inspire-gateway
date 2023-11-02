package scope

import (
	"context"

	scopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/scope"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"
	scopemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/scope"
	"github.com/google/uuid"
)

func (h *Handler) CreateScope(ctx context.Context) (*npool.Scope, error) {
	id := uuid.NewString()
	if h.ID == nil {
		h.ID = &id
	}
	if _, err := scopemwcli.CreateScope(
		ctx,
		&scopemwpb.ScopeReq{
			ID:          h.ID,
			GoodID:      h.GoodID,
			CouponID:    h.CouponID,
			CouponScope: h.CouponScope,
		},
	); err != nil {
		return nil, err
	}

	return h.GetScope(ctx)
}
