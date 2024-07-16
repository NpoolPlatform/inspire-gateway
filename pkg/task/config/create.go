package config

import (
	"context"

	"github.com/google/uuid"

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
		EntID:            h.EntID,
		AppID:            h.AppID,
		EventID:          h.EventID,
		TaskType:         h.TaskType,
		Name:             h.Name,
		TaskDesc:         h.TaskDesc,
		StepGuide:        h.StepGuide,
		RecommendMessage: h.RecommendMessage,
		Index:            h.Index,
		MaxRewardCount:   h.MaxRewardCount,
		CooldownSecord:   h.CooldownSecord,
		LastTaskID:       h.LastTaskID,
	}); err != nil {
		return nil, err
	}

	return h.GetTaskConfig(ctx, nil)
}
