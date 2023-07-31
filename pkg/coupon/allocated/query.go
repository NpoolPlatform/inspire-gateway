package allocated

import (
	"context"
	"fmt"

	allocatedmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/allocated"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	allocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/allocated"
)

func (h *Handler) GetCoupons(ctx context.Context) ([]*allocatedmwpb.Coupon, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appid")
	}

	conds := &allocatedmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}

	return allocatedmwcli.GetCoupons(
		ctx,
		conds,
		h.Offset,
		h.Limit,
	)
}
