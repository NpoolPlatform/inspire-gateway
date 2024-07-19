package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	configmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"
	configmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/config"
)

func (h *Handler) UpdateTaskConfig(ctx context.Context) (*npool.TaskConfig, error) {
	info, err := h.GetTaskConfig(ctx, nil)
	if err != nil {
		return nil, wlog.WrapError(err)
	}
	if info == nil {
		return nil, wlog.Errorf("invalid config")
	}
	if info.ID != *h.ID || info.EntID != *h.EntID {
		return nil, wlog.Errorf("permission denied")
	}

	if err := configmwcli.UpdateTaskConfig(ctx, &configmwpb.TaskConfigReq{
		ID:               h.ID,
		EntID:            h.EntID,
		AppID:            h.AppID,
		EventID:          h.EventID,
		TaskType:         h.TaskType,
		Name:             h.Name,
		TaskDesc:         h.TaskDesc,
		StepGuide:        h.StepGuide,
		RecommendMessage: h.RecommendMessage,
		Index:            h.Index,
		LastTaskID:       h.LastTaskID,
		MaxRewardCount:   h.MaxRewardCount,
		CooldownSecord:   h.CooldownSecord,
	}); err != nil {
		return nil, wlog.WrapError(err)
	}
	return h.GetTaskConfig(ctx, nil)
}
