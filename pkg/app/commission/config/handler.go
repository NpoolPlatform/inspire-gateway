package config

import (
	"context"
	"fmt"
	"time"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Handler struct {
	ID              *uint32
	EntID           *string
	AppID           *string
	SettleType      *types.SettleType
	AmountOrPercent *string
	ThresholdAmount *string
	Invites         *uint32
	StartAt         *uint32
	EndAt           *uint32
	CheckAffiliate  bool
	Offset          int32
	Limit           int32
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	handler := &Handler{
		CheckAffiliate: true,
	}
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
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.AppID = id
		return nil
	}
}

func WithSettleType(settleType *types.SettleType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if settleType == nil {
			if must {
				return fmt.Errorf("invalid settletype")
			}
			return nil
		}
		switch *settleType {
		case types.SettleType_GoodOrderPayment:
		case types.SettleType_TechniqueServiceFee:
		default:
			return fmt.Errorf("invalid settletype")
		}
		h.SettleType = settleType
		return nil
	}
}

func WithAmountOrPercent(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid amountorpercent")
			}
			return nil
		}
		if _, err := decimal.NewFromString(*amount); err != nil {
			return err
		}
		h.AmountOrPercent = amount
		return nil
	}
}

func WithThresholdAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid thresholdamount")
			}
			return nil
		}
		if _, err := decimal.NewFromString(*amount); err != nil {
			return err
		}
		h.ThresholdAmount = amount
		return nil
	}
}

func WithInvites(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid invites")
			}
			return nil
		}
		h.Invites = value
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
		h.EndAt = value
		return nil
	}
}

func WithCheckAffiliate(check bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.CheckAffiliate = check
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