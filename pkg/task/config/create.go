package config

import (
	"context"

	"github.com/google/uuid"

	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	taskconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"
	taskconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/config"
)

func (h *Handler) CreateTaskConfig(ctx context.Context) (*npool.TaskConfig, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	if err := taskconfigmwcli.CreateTaskConfig(ctx, &taskconfigmwpb.TaskConfigReq{
		EntID:                  h.EntID,
		AppID:                  h.AppID,
		EventID:                h.EventID,
		TaskType:               h.TaskType,
		Name:                   h.Name,
		TaskDesc:               h.TaskDesc,
		StepGuide:              h.StepGuide,
		RecommendMessage:       h.RecommendMessage,
		Index:                  h.Index,
		MaxRewardCount:         h.MaxRewardCount,
		CooldownSecond:         h.CooldownSecond,
		LastTaskID:             h.LastTaskID,
		IntervalReset:          h.IntervalReset,
		IntervalResetSecond:    h.IntervalResetSecond,
		MaxIntervalRewardCount: h.MaxIntervalRewardCount,
	}); err != nil {
		return nil, wlog.WrapError(err)
	}

	return h.GetTaskConfig(ctx, nil)
}
