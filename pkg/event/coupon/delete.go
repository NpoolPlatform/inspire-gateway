package coupon

import (
	"context"
	"fmt"

	eventcouponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event/coupon"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coupon"
	eventcouponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event/coupon"
)

func (h *Handler) DeleteEventCoupon(ctx context.Context) (*npool.EventCoupon, error) {
	info, err := eventcouponmwcli.GetEventCouponOnly(ctx, &eventcouponmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid eventcoupon")
	}

	if err := eventcouponmwcli.DeleteEventCoupon(ctx, h.ID, h.EntID); err != nil {
		return nil, err
	}

	return h.GetEventCouponExt(ctx, info)
}
