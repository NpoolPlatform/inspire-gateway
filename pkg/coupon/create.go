//nolint:dupl
package coupon

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	couponmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon"
)

type createHandler struct {
	*Handler
	goodID *string
}

func (h *createHandler) checkGood(ctx context.Context) error {
	if h.AppGoodID == nil {
		return nil
	}
	info, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid good")
	}
	h.goodID = &info.GoodID
	return nil
}

func (h *createHandler) validateIssuer(ctx context.Context) error {
	if h.IssuedBy == nil {
		return fmt.Errorf("invalid issuer")
	}
	exist, err := usermwcli.ExistUserConds(ctx, &usermwpb.Conds{
		ID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.IssuedBy},
	})
	if err != nil {
		return nil
	}
	if !exist {
		return fmt.Errorf("invalid issuer")
	}
	return nil
}

func (h *createHandler) createSpecialOffer(ctx context.Context) (*couponmwpb.Coupon, error) {
	// TODO: need dtm to create coupon and allocated
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}

func (h *Handler) CreateCoupon(ctx context.Context) (*couponmwpb.Coupon, error) {
	if h.CouponType == nil {
		return nil, fmt.Errorf("invalid coupontype")
	}
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appid")
	}

	handler := &createHandler{
		Handler: h,
	}

	if err := handler.checkGood(ctx); err != nil {
		return nil, err
	}
	if err := handler.validateIssuer(ctx); err != nil {
		return nil, err
	}

	if *h.CouponType == types.CouponType_SpecialOffer {
		return handler.createSpecialOffer(ctx)
	}
	return couponmwcli.CreateCoupon(ctx, &couponmwpb.CouponReq{
		ID:               h.ID,
		AppID:            h.AppID,
		CouponType:       h.CouponType,
		Denomination:     h.Denomination,
		Circulation:      h.Circulation,
		IssuedBy:         h.IssuedBy,
		StartAt:          h.StartAt,
		DurationDays:     h.DurationDays,
		Message:          h.Message,
		Name:             h.Name,
		GoodID:           handler.goodID,
		AppGoodID:        h.AppGoodID,
		CouponConstraint: h.CouponConstraint,
		Threshold:        h.Threshold,
		Random:           h.Random,
		CouponScope:      h.CouponScope,
	})
}
