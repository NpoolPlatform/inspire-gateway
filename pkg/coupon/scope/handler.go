package scope

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	couponmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	"github.com/google/uuid"
)

type Handler struct {
	ID          *string
	AppID       *string
	UserID      *string
	AppGoodID   *string
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

func WithID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid id")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ID = id
		return nil
	}
}

//nolint
func WithAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid appid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}

		exist, err := appmwcli.ExistApp(ctx, *id)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("app not found")
		}

		h.AppID = id
		return nil
	}
}

func WithUserID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid userid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.UserID = id
		return nil
	}
}

//nolint
func WithAppGoodID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid appgoodid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		exist, err := appgoodmwcli.ExistGood(ctx, *id)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("appgood not found")
		}
		h.AppGoodID = id
		return nil
	}
}

//nolint
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
