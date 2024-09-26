package event

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Handler struct {
	ID             *uint32
	EntID          *string
	AppID          *string
	EventType      *basetypes.UsedFor
	Credits        *string
	CreditsPerUSD  *string
	MaxConsecutive *uint32
	AppGoodID      *string
	InviterLayers  *uint32
	Offset         int32
	Limit          int32
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

//nolint:gocyclo
func WithEventType(eventType *basetypes.UsedFor, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if eventType == nil {
			if must {
				return fmt.Errorf("invalid eventtype")
			}
			return nil
		}
		switch *eventType {
		case basetypes.UsedFor_Signup:
		case basetypes.UsedFor_Signin:
		case basetypes.UsedFor_Update:
		case basetypes.UsedFor_SetWithdrawAddress:
		case basetypes.UsedFor_Withdraw:
		case basetypes.UsedFor_CreateInvitationCode:
		case basetypes.UsedFor_SetCommission:
		case basetypes.UsedFor_SetTransferTargetUser:
		case basetypes.UsedFor_Transfer:
		case basetypes.UsedFor_WithdrawalRequest:
		case basetypes.UsedFor_Purchase:
		case basetypes.UsedFor_SimulateOrderProfit:
		case basetypes.UsedFor_ConsecutiveLogin:
		case basetypes.UsedFor_GoodSocialSharing:
		case basetypes.UsedFor_FirstOrderCompleted:
		case basetypes.UsedFor_SetAddress:
		case basetypes.UsedFor_Set2FA:
		case basetypes.UsedFor_FirstBenefit:
		case basetypes.UsedFor_WriteComment:
		case basetypes.UsedFor_WriteRecommend:
		case basetypes.UsedFor_GoodScoring:
		case basetypes.UsedFor_SubmitTicket:
		case basetypes.UsedFor_IntallApp:
		case basetypes.UsedFor_SetNFTAvatar:
		case basetypes.UsedFor_SetPersonalImage:
		case basetypes.UsedFor_KYCApproved:
		case basetypes.UsedFor_OrderCompleted:
		case basetypes.UsedFor_WithdrawalCompleted:
		case basetypes.UsedFor_DepositReceived:
		case basetypes.UsedFor_UpdatePassword:
		case basetypes.UsedFor_ResetPassword:
		case basetypes.UsedFor_InternalTransfer:
		case basetypes.UsedFor_Contact:
		case basetypes.UsedFor_KYCRejected:
		case basetypes.UsedFor_NewLogin:
		}
		h.EventType = eventType
		return nil
	}
}

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
		h.AppGoodID = id
		return nil
	}
}

func WithCredits(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid amount")
			}
			return nil
		}
		if _, err := decimal.NewFromString(*amount); err != nil {
			return err
		}
		h.Credits = amount
		return nil
	}
}

func WithCreditsPerUSD(value *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid creditsperusd")
			}
			return nil
		}
		if _, err := decimal.NewFromString(*value); err != nil {
			return err
		}
		h.CreditsPerUSD = value
		return nil
	}
}

func WithMaxConsecutive(consecutive *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if consecutive == nil {
			if must {
				return fmt.Errorf("invalid consecutive")
			}
			return nil
		}
		h.MaxConsecutive = consecutive
		return nil
	}
}

func WithInviterLayers(layers *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if layers == nil {
			if must {
				return fmt.Errorf("invalid layers")
			}
			return nil
		}
		h.InviterLayers = layers
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
