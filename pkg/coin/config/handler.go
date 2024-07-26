package config

import (
	"context"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	"github.com/shopspring/decimal"

	"github.com/google/uuid"
)

type Handler struct {
	ID         *uint32
	EntID      *string
	AppID      *string
	CoinTypeID *string
	MaxValue   *string
	Allocated  *string
	Offset     int32
	Limit      int32
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
				return wlog.Errorf("invalid id")
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
				return wlog.Errorf("invalid entid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return wlog.WrapError(err)
		}
		h.EntID = id
		return nil
	}
}

func WithAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return wlog.Errorf("invalid appid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return wlog.WrapError(err)
		}
		exist, err := appmwcli.ExistApp(ctx, *id)
		if err != nil {
			return wlog.WrapError(err)
		}
		if !exist {
			return wlog.Errorf("invalid appid")
		}
		h.AppID = id
		return nil
	}
}

func WithCoinTypeID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return wlog.Errorf("invalid cointypeid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return wlog.WrapError(err)
		}
		h.CoinTypeID = id
		return nil
	}
}

//nolint:dupl
func WithMaxValue(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return wlog.Errorf("invalid maxvalue")
			}
			return nil
		}
		_amount, err := decimal.NewFromString(*amount)
		if err != nil {
			return wlog.WrapError(err)
		}
		if _amount.Cmp(decimal.NewFromInt(0)) < 0 {
			return wlog.Errorf("invalid maxvalue")
		}
		h.MaxValue = amount
		return nil
	}
}

//nolint:dupl
func WithAllocated(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return wlog.Errorf("invalid allocated")
			}
			return nil
		}
		_amount, err := decimal.NewFromString(*amount)
		if err != nil {
			return wlog.WrapError(err)
		}
		if _amount.Cmp(decimal.NewFromInt(0)) < 0 {
			return wlog.Errorf("invalid allocated")
		}
		h.Allocated = amount
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
