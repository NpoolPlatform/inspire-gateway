package config

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID               *uint32
	EntID            *string
	AppID            *string
	SettleMode       *types.SettleMode
	SettleAmountType *types.SettleAmountType
	SettleInterval   *types.SettleInterval
	CommissionType   *types.CommissionType
	SettleBenefit    *bool
	StartAt          *uint32
	EndAt            *uint32
	MaxLevel         *uint32
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

func WithCommissionType(commissionType *types.CommissionType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if commissionType == nil {
			if must {
				return fmt.Errorf("invalid commissiontype")
			}
			return nil
		}
		switch *commissionType {
		case types.CommissionType_LegacyCommission:
		case types.CommissionType_LayeredCommission:
		case types.CommissionType_DirectCommission:
		default:
			return fmt.Errorf("invalid commissiontype")
		}
		h.CommissionType = commissionType
		return nil
	}
}

func WithSettleAmountType(settleAmountType *types.SettleAmountType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if settleAmountType == nil {
			if must {
				return fmt.Errorf("invalid settleamounttype")
			}
			return nil
		}
		switch *settleAmountType {
		case types.SettleAmountType_SettleByPercent:
		case types.SettleAmountType_SettleByAmount:
		default:
			return fmt.Errorf("invalid settleamounttype")
		}
		h.SettleAmountType = settleAmountType
		return nil
	}
}

func WithSettleMode(settleMode *types.SettleMode, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if settleMode == nil {
			if must {
				return fmt.Errorf("invalid settlemode")
			}
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

func WithSettleInterval(settleInterval *types.SettleInterval, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if settleInterval == nil {
			if must {
				return fmt.Errorf("invalid settleinterval")
			}
			return nil
		}
		switch *settleInterval {
		case types.SettleInterval_SettleAggregate:
		case types.SettleInterval_SettleYearly:
		case types.SettleInterval_SettleMonthly:
		case types.SettleInterval_SettleEveryOrder:
		default:
			return fmt.Errorf("invalid settleinterval")
		}
		h.SettleInterval = settleInterval
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
		h.StartAt = value
		return nil
	}
}

func WithEndAt(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.EndAt = value
		return nil
	}
}

func WithSettleBenefit(value *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid settlebenefit")
			}
			return nil
		}
		h.SettleBenefit = value
		return nil
	}
}

func WithMaxLevel(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid MaxLevel")
			}
			return nil
		}
		h.MaxLevel = value
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
