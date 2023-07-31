package coupon

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Handler struct {
	ID               *string
	AppID            *string
	UserID           *string
	IssuedBy         *string
	GoodID           *string
	CouponType       *types.CouponType
	Denomination     *string
	Circulation      *string
	StartAt          *uint32
	DurationDays     *uint32
	Message          *string
	Name             *string
	Threshold        *string
	CouponConstraint *types.CouponConstraint
	Random           *bool
	Offset           int32
	Limit            int32
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

func WithID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ID = id
		return nil
	}
}

func WithAppID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.UserID = id
		return nil
	}
}

func WithIssuedBy(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.IssuedBy = id
		return nil
	}
}

func WithGoodID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.GoodID = id
		return nil
	}
}

func WithCouponType(couponType *types.CouponType) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if couponType == nil {
			return nil
		}
		switch *couponType {
		case types.CouponType_FixAmount:
		case types.CouponType_Discount:
		case types.CouponType_SpecialOffer:
		default:
			return fmt.Errorf("invalid coupontype")
		}
		h.CouponType = couponType
		return nil
	}
}

func WithDenomination(amount *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			return nil
		}
		if _, err := decimal.NewFromString(*amount); err != nil {
			return err
		}
		h.Denomination = amount
		return nil
	}
}

func WithCirculation(amount *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			return nil
		}
		if _, err := decimal.NewFromString(*amount); err != nil {
			return err
		}
		h.Circulation = amount
		return nil
	}
}

func WithStartAt(value *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.StartAt = value
		return nil
	}
}

func WithDurationDays(value *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.DurationDays = value
		return nil
	}
}

func WithMessage(value *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Message = value
		return nil
	}
}

func WithName(value *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			return nil
		}
		const leastNameLen = 3
		if len(*value) < leastNameLen {
			return fmt.Errorf("invalid name")
		}
		h.Name = value
		return nil
	}
}

func WithThreshold(amount *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			return nil
		}
		if _, err := decimal.NewFromString(*amount); err != nil {
			return err
		}
		h.Threshold = amount
		return nil
	}
}

func WithCouponConstraint(couponConstraint *types.CouponConstraint) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if couponConstraint == nil {
			return nil
		}
		switch *couponConstraint {
		case types.CouponConstraint_Normal:
		case types.CouponConstraint_PaymentThreshold:
		case types.CouponConstraint_GoodOnly:
		case types.CouponConstraint_GoodThreshold:
		default:
			return fmt.Errorf("invalid couponconstraint")
		}
		h.CouponConstraint = couponConstraint
		return nil
	}
}

func WithRandom(value *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Random = value
		return nil
	}
}

func WithOffset(offset int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Offset = offset
		return nil
	}
}

func WithLimit(limit int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if limit == 0 {
			limit = constant.DefaultRowLimit
		}
		h.Limit = limit
		return nil
	}
}
