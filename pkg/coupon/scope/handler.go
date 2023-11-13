package scope

import (
	"context"
	"fmt"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	"github.com/google/uuid"
)

type Handler struct {
	ID          *uint32
	EntID       *string
	GoodID      *string
	CouponID    *string
	CouponScope *types.CouponScope
	Offset      int32
	Limit       int32
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	handler := &Handler{}
	for _, opt := range options {
		if err := opt(ctx, handler); err != nil {
			return nil, err
		}
	}
	return handler, nil
}

func WithID(id *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid id")
			}
			return nil
		}
		h.ID = id
		return nil
	}
}

func WithEntID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid entid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.EntID = id
		return nil
	}
}

func WithGoodID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid goodid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		good, err := goodmwcli.GetGood(ctx, *id)
		if err != nil {
			return err
		}
		if good == nil {
			return fmt.Errorf("good not found")
		}
		h.GoodID = id
		return nil
	}
}

func WithCouponID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid couponid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		exist, err := couponmwcli.ExistCoupon(ctx, *id)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("coupon not found")
		}
		h.CouponID = id
		return nil
	}
}

func WithCouponScope(couponScope *types.CouponScope, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if couponScope == nil {
			if must {
				return fmt.Errorf("invalid couponscope")
			}
			return nil
		}
		switch *couponScope {
		case types.CouponScope_Blacklist:
		case types.CouponScope_Whitelist:
		default:
			return fmt.Errorf("invalid couponscope")
		}
		h.CouponScope = couponScope
		return nil
	}
}

func WithOffset(value int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Offset = value
		return nil
	}
}

func WithLimit(value int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == 0 {
			value = constant.DefaultRowLimit
		}
		h.Limit = value
		return nil
	}
}
