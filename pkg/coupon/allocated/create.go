package allocated

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	allocatedmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/allocated"
	allocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/allocated"

	"github.com/google/uuid"
)

func (h *Handler) CreateCoupon(ctx context.Context) (*allocatedmwpb.Coupon, error) {
	exist, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid user")
	}

	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	if _, err := allocatedmwcli.CreateCoupon(
		ctx,
		&allocatedmwpb.CouponReq{
			EntID:    h.EntID,
			AppID:    h.AppID,
			UserID:   h.UserID,
			CouponID: h.CouponID,
		},
	); err != nil {
		return nil, err
	}

	return h.GetCoupon(ctx)
}
