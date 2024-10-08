package config

import (
	"context"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID                     *uint32
	EntID                  *string
	AppID                  *string
	EventID                *string
	TaskType               *types.TaskType
	Name                   *string
	TaskDesc               *string
	StepGuide              *string
	RecommendMessage       *string
	Index                  *uint32
	LastTaskID             *string
	UserID                 *string
	MaxRewardCount         *uint32
	CooldownSecond         *uint32
	IntervalReset          *bool
	IntervalResetSecond    *uint32
	MaxIntervalRewardCount *uint32
	Offset                 int32
	Limit                  int32
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

func WithEventID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return wlog.Errorf("invalid eventid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return wlog.WrapError(err)
		}
		h.EventID = id
		return nil
	}
}

func WithTaskType(taskType *types.TaskType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if taskType == nil {
			if must {
				return wlog.Errorf("invalid tasktype")
			}
			return nil
		}
		switch *taskType {
		case types.TaskType_BaseTask:
		case types.TaskType_GrowthTask:
		default:
			return wlog.Errorf("invalid tasktype")
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
				return wlog.Errorf("invalid index")
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
				return wlog.Errorf("invalid lasttaskid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return wlog.WrapError(err)
		}
		h.LastTaskID = id
		return nil
	}
}

func WithUserID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return wlog.Errorf("invalid userid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return wlog.WrapError(err)
		}
		h.UserID = id
		return nil
	}
}

func WithMaxRewardCount(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return wlog.Errorf("invalid maxrewardcount")
			}
			return nil
		}
		h.MaxRewardCount = value
		return nil
	}
}

func WithCooldownSecond(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return wlog.Errorf("invalid cooldownsecord")
			}
			return nil
		}
		h.CooldownSecond = value
		return nil
	}
}

func WithIntervalReset(value *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return wlog.Errorf("invalid intervalreset")
			}
			return nil
		}
		h.IntervalReset = value
		return nil
	}
}

func WithIntervalResetSecond(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return wlog.Errorf("invalid intervalresetsecond")
			}
			return nil
		}
		h.IntervalResetSecond = value
		return nil
	}
}

func WithMaxIntervalRewardCount(value *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if value == nil {
			if must {
				return wlog.Errorf("invalid maxintervalrewardcount")
			}
			return nil
		}
		h.MaxIntervalRewardCount = value
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
