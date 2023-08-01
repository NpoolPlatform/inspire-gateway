package commission

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Handler struct {
	ID              *string
	AppID           *string
	UserID          *string
	GoodID          *string
	SettleType      *types.SettleType
	SettleMode      *types.SettleMode
	AmountOrPercent *string
	StartAt         *uint32
	EndAt           *uint32
	FromGoodID      *string
	ToGoodID        *string
	ScalePercent    *string
	Offset          int32
	Limit           int32
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

func WithUserID(id *string) func(context.Context, *Handler) error {
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

func WithSettleType(settleType *types.SettleType) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if settleType == nil {
			return nil
		}
		switch *settleType {
		case types.SettleType_GoodOrderPercent:
		case types.SettleType_TechniqueFeePercent:
		case types.SettleType_NoCommission:
		default:
			return fmt.Errorf("invalid settletype")
		}
		h.SettleType = settleType
		return nil
	}
}

func WithSettleMode(settleMode *types.SettleMode) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if settleMode == nil {
			return nil
		}
		switch *settleMode {
		case types.SettleMode_SettleWithPaymentAmount:
		case types.SettleMode_SettleWithGoodValue:
		default:
			return fmt.Errorf("invalid settlemode")
		}
		h.SettleMode = settleMode
		return nil
	}
}

func WithAmountOrPercent(amount *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			return nil
		}
		if _, err := decimal.NewFromString(*amount); err != nil {
			return err
		}
		h.AmountOrPercent = amount
		return nil
	}
}

func WithStartAt(value *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.StartAt = value
		return nil
	}
}

func WithEndAt(value *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.EndAt = value
		return nil
	}
}

func WithFromGoodID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.FromGoodID = id
		return nil
	}
}

func WithToGoodID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ToGoodID = id
		return nil
	}
}

func WithScalePercent(percent *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if percent == nil {
			return nil
		}
		if _, err := decimal.NewFromString(*percent); err != nil {
			return err
		}
		h.ScalePercent = percent
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