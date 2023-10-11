package scope

import (
	"context"
	"fmt"

	scopemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/scope"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"
	scopemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/scope"
	"github.com/google/uuid"
)

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
		return nil, err
	}
	return h.GetScope(ctx)
}
