package coupon

import (
	"context"
	"fmt"
	"time"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	timedef "github.com/NpoolPlatform/go-service-framework/pkg/const/time"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Handler struct {
	ID                  *uint32
	EntID               *string
	AppID               *string
	TargetAppID         *string
	UserID              *string
	IssuedBy            *string
	CouponType          *types.CouponType
	Denomination        *string
	Circulation         *string
	StartAt             *uint32
	EndAt               *uint32
	DurationDays        *uint32
	Message             *string
	Name                *string
	Threshold           *string
	CouponConstraint    *types.CouponConstraint
	CouponScope         *types.CouponScope
	Random              *bool
	CashableProbability *string
	Offset              int32
	Limit               int32
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

func WithAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid appid")
			}
			return nil
		}
		exist, err := appmwcli.ExistApp(ctx, *id)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("invalid appid")
		}
		h.AppID = id
		return nil
	}
}

func WithTargetAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid targetappid")
			}
			return nil
		}
		exist, err := appmwcli.ExistApp(ctx, *id)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("invalid targetappid")
		}
		h.TargetAppID = id
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

func WithIssuedBy(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid issueby")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.IssuedBy = id
		return nil
	}
}

func WithCouponType(couponType *types.CouponType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if couponType == nil {
			if must {
				return fmt.Errorf("invalid coupontype")
			}
			return nil
		}
		switch *couponType {
		case types.CouponType_FixAmount:
		case types.CouponType_Discount:
		default:
			return fmt.Errorf("invalid coupontype")
		}
		h.CouponType = couponType
		return nil
	}
}

func WithDenomination(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid denomination")
			}
			return nil
		}
		if _, err := decimal.NewFromString(*amount); err != nil {
			return err
		}
		h.Denomination = amount
		return nil
	}
}

func WithCirculation(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid circulation")
			}
			return nil
		}
		if _, err := decimal.NewFromString(*amount); err != nil {
			return err
		}
		h.Circulation = amount
		return nil
	}
}

func WithStartAt(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid startat")
			}
			return nil
		}
		if *value == 0 {
			*value = uint32(time.Now().Unix())
		}
		h.StartAt = value
		return nil
	}
}

func WithEndAt(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid endat")
			}
			return nil
		}
		if *value == 0 {
			return fmt.Errorf("invalid endat")
		}
		h.EndAt = value
		return nil
	}
}
func WithDurationDays(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid duration days")
			}
			return nil
		}
		if *value == 0 {
			*value = timedef.DaysPerYear
		}
		h.DurationDays = value
		return nil
	}
}

func WithMessage(value *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Message = value
		return nil
	}
}

func WithName(value *string, must bool) func(context.Context, *Handler) error {
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

func WithThreshold(amount *string, must bool) func(context.Context, *Handler) error {
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

func WithCashableProbability(probability *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if probability == nil {
			if must {
				return fmt.Errorf("invalid probability")
			}
			return nil
		}
		_probability, err := decimal.NewFromString(*probability)
		if err != nil {
			return err
		}
		if _probability.Cmp(decimal.NewFromInt(0)) < 0 || _probability.Cmp(decimal.NewFromInt(1)) > 0 {
			return fmt.Errorf("invalid probability")
		}
		h.CashableProbability = probability
		return nil
	}
}

func WithCouponConstraint(couponConstraint *types.CouponConstraint, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if couponConstraint == nil {
			return nil
		}
		switch *couponConstraint {
		case types.CouponConstraint_Normal:
		case types.CouponConstraint_PaymentThreshold:
		default:
			return fmt.Errorf("invalid couponconstraint")
		}
		h.CouponConstraint = couponConstraint
		return nil
	}
}

func WithCouponScope(couponScope *types.CouponScope, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if couponScope == nil {
			return nil
		}
		switch *couponScope {
		case types.CouponScope_AllGood:
		case types.CouponScope_Blacklist:
		case types.CouponScope_Whitelist:
		default:
			return fmt.Errorf("invalid couponscope")
		}
		h.CouponScope = couponScope
		return nil
	}
}

func WithRandom(value *bool, must bool) func(context.Context, *Handler) error {
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
