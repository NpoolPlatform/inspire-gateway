package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	taskconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	taskconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/config"
)

type queryHandler struct {
	*Handler
	taskConfigs []*taskconfigmwpb.TaskConfig
	events      map[string]*eventmwpb.Event
	taskInfos   []*npool.TaskConfig
}

func (h *queryHandler) formalize() {
	for _, val := range h.taskConfigs {
		task := &npool.TaskConfig{
			ID:               val.ID,
			EntID:            val.EntID,
			AppID:            val.AppID,
			EventID:          val.EventID,
			TaskType:         val.TaskType,
			Name:             val.Name,
			TaskDesc:         val.TaskDesc,
			StepGuide:        val.StepGuide,
			RecommendMessage: val.RecommendMessage,
			Index:            val.Index,
			LastTaskID:       val.LastTaskID,
			MaxRewardCount:   val.MaxRewardCount,
			CooldownSecond:   val.CooldownSecond,
			CreatedAt:        val.CreatedAt,
			UpdatedAt:        val.UpdatedAt,
		}
		event, ok := h.events[val.EventID]
		if ok {
			task.EventType = event.EventType
		}
		h.taskInfos = append(h.taskInfos, task)
	}
}

func (h *queryHandler) getEvents(ctx context.Context) error {
	eventIDs := []string{}
	for _, val := range h.taskConfigs {
		eventIDs = append(eventIDs, val.EventID)
	}
	events, _, err := eventmwcli.GetEvents(ctx, &eventmwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: eventIDs},
	}, 0, int32(len(eventIDs)))
	if err != nil {
		return wlog.WrapError(err)
	}

	for _, event := range events {
		h.events[event.EntID] = event
	}
	return nil
}

func (h *Handler) GetTaskConfig(ctx context.Context, info *taskconfigmwpb.TaskConfig) (*npool.TaskConfig, error) {
	if h.EntID == nil {
		return nil, wlog.Errorf("invalid entid")
	}

	if info == nil {
		taskInfo, err := taskconfigmwcli.GetTaskConfig(ctx, *h.EntID)
		if err != nil {
			return nil, wlog.WrapError(err)
		}
		if taskInfo == nil {
			return nil, nil
		}
		info = taskInfo
	}

	handler := &queryHandler{
		Handler:     h,
		taskConfigs: []*taskconfigmwpb.TaskConfig{info},
		events:      map[string]*eventmwpb.Event{},
		taskInfos:   []*npool.TaskConfig{},
	}

	if err := handler.getEvents(ctx); err != nil {
		return nil, wlog.WrapError(err)
	}

	handler.formalize()
	if len(handler.taskInfos) == 0 {
		return nil, nil
	}
	return handler.taskInfos[0], nil
}

func (h *Handler) GetTaskConfigs(ctx context.Context) ([]*npool.TaskConfig, uint32, error) {
	handler := &queryHandler{
		Handler:     h,
		taskConfigs: []*taskconfigmwpb.TaskConfig{},
		events:      map[string]*eventmwpb.Event{},
		taskInfos:   []*npool.TaskConfig{},
	}

	conds := &taskconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}

	infos, total, err := taskconfigmwcli.GetTaskConfigs(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, wlog.WrapError(err)
	}
	if len(infos) == 0 {
		return nil, total, nil
	}
	handler.taskConfigs = infos

	if err := handler.getEvents(ctx); err != nil {
		return nil, 0, wlog.WrapError(err)
	}

	handler.formalize()

	return handler.taskInfos, total, nil
}
