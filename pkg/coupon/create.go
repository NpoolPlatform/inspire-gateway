package coupon

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) validateIssuer(ctx context.Context) error {
	if h.IssuedBy == nil {
		return fmt.Errorf("invalid issuer")
	}
	exist, err := usermwcli.ExistUserConds(ctx, &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.IssuedBy},
	})
	if err != nil {
		return nil
	}
	if !exist {
		return fmt.Errorf("invalid issuer")
	}
	return nil
}

func (h *Handler) CreateCoupon(ctx context.Context) (*couponmwpb.Coupon, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	handler := &createHandler{
		Handler: h,
	}
	if err := handler.validateIssuer(ctx); err != nil {
		return nil, err
	}
	return couponmwcli.CreateCoupon(ctx, &couponmwpb.CouponReq{
		EntID:                         h.EntID,
		AppID:                         h.AppID,
		CouponType:                    h.CouponType,
		Denomination:                  h.Denomination,
		Circulation:                   h.Circulation,
		IssuedBy:                      h.IssuedBy,
		StartAt:                       h.StartAt,
		EndAt:                         h.EndAt,
		DurationDays:                  h.DurationDays,
		Message:                       h.Message,
		Name:                          h.Name,
		CouponConstraint:              h.CouponConstraint,
		Threshold:                     h.Threshold,
		Random:                        h.Random,
		CouponScope:                   h.CouponScope,
		CashableProbabilityPerMillion: h.CashableProbabilityPerMillion,
	})
}
