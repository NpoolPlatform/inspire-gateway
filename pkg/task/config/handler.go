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
	EventID          *string
	TaskType         *types.TaskType
	Name             *string
	TaskDesc         *string
	StepGuide        *string
	RecommendMessage *string
	Index            *uint32
	LastTaskID       *string
	UserID           *string
	MaxRewardCount   *uint32
	CooldownSecord   *uint32
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

func WithEventID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid eventid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.EventID = id
		return nil
	}
}

func WithTaskType(taskType *types.TaskType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if taskType == nil {
			if must {
				return fmt.Errorf("invalid tasktype")
			}
			return nil
		}
		switch *taskType {
		case types.TaskType_BaseTask:
		case types.TaskType_DailyTask:
		case types.TaskType_GrowthTask:
		default:
			return fmt.Errorf("invalid tasktype")
		}
		h.TaskType = taskType
		return nil
	}
}

func WithName(value *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Name = value
		return nil
	}
}

func WithTaskDesc(value *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.TaskDesc = value
		return nil
	}
}

func WithStepGuide(value *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.StepGuide = value
		return nil
	}
}

func WithRecommendMessage(value *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.RecommendMessage = value
		return nil
	}
}

func WithIndex(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid index")
			}
			return nil
		}
		h.Index = value
		return nil
	}
}

func WithLastTaskID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid lasttaskid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.LastTaskID = id
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

func WithMaxRewardCount(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid maxrewardcount")
			}
			return nil
		}
		h.MaxRewardCount = value
		return nil
	}
}

func WithCooldownSecord(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return fmt.Errorf("invalid cooldownsecord")
			}
			return nil
		}
		h.CooldownSecord = value
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
