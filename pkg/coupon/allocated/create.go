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
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appid")
	}
	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userid")
	}

	exist, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid user")
	}

	id := uuid.NewString()
	if h.ID == nil {
		h.ID = &id
	}

	if _, err := allocatedmwcli.CreateCoupon(
		ctx,
		&allocatedmwpb.CouponReq{
			ID:       h.ID,
			AppID:    h.AppID,
			UserID:   h.UserID,
			CouponID: h.CouponID,
		},
	); err != nil {
		return nil, err
	}

	return h.GetCoupon(ctx)
}
