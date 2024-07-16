package config

import (
	"context"
	"fmt"

	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	taskconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/config"
	taskusermwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/user"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	taskconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/config"
	taskusermwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/user"
)

type queryHandler struct {
	*Handler
	taskConfigs   []*taskconfigmwpb.TaskConfig
	taskUsers     []*taskusermwpb.TaskUser
	events        map[string]*eventmwpb.Event
	userTaskInfos []*npool.UserTaskConfig
	taskInfos     []*npool.TaskConfig
}

func (h *queryHandler) formalizeUserTask() {
	taskMap := map[string]uint32{}
	for _, taskUser := range h.taskUsers {
		taskUserCount, ok := taskMap[taskUser.TaskID]
		if ok {
			taskUserCount++
			taskMap[taskUser.TaskID] = taskUserCount
			continue
		}
		taskMap[taskUser.TaskID] = uint32(1)
	}
	for _, comm := range h.taskConfigs {
		taskState := types.TaskState_NotStarted
		rewardState := types.RewardState_UnIssued
		taskUserCount, ok := taskMap[comm.EntID]
		if ok {
			taskState = types.TaskState_Done
			rewardState = types.RewardState_Issued
		}
		h.userTaskInfos = append(h.userTaskInfos, &npool.UserTaskConfig{
			ID:               comm.ID,
			EntID:            comm.EntID,
			AppID:            comm.AppID,
			EventID:          comm.EventID,
			TaskType:         comm.TaskType,
			Name:             comm.Name,
			TaskDesc:         comm.TaskDesc,
			StepGuide:        comm.StepGuide,
			RecommendMessage: comm.RecommendMessage,
			Index:            comm.Index,
			LastTaskID:       comm.LastTaskID,
			MaxRewardCount:   comm.MaxRewardCount,
			CooldownSecord:   comm.CooldownSecord,
			CompletionTimes:  taskUserCount,
			TaskState:        taskState,
			RewardState:      rewardState,
			CreatedAt:        comm.CreatedAt,
			UpdatedAt:        comm.UpdatedAt,
		})
	}
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
			CooldownSecord:   val.CooldownSecord,
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

func (h *Handler) GetUserTaskConfig(ctx context.Context) (*npool.UserTaskConfig, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}

	info, err := taskconfigmwcli.GetTaskConfig(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:       h,
		taskConfigs:   []*taskconfigmwpb.TaskConfig{info},
		taskUsers:     []*taskusermwpb.TaskUser{},
		userTaskInfos: []*npool.UserTaskConfig{},
	}

	handler.formalizeUserTask()
	if len(handler.userTaskInfos) == 0 {
		return nil, nil
	}

	return handler.userTaskInfos[0], nil
}

func (h *Handler) GetUserTaskConfigs(ctx context.Context) ([]*npool.UserTaskConfig, uint32, error) {
	handler := &queryHandler{
		Handler:       h,
		userTaskInfos: []*npool.UserTaskConfig{},
		taskConfigs:   []*taskconfigmwpb.TaskConfig{},
		taskUsers:     []*taskusermwpb.TaskUser{},
	}

	conds := &taskconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}

	infos, total, err := taskconfigmwcli.GetTaskConfigs(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}
	handler.taskConfigs = infos

	userConds := &taskusermwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
	}

	userInfos, _, err := taskusermwcli.GetTaskUsers(ctx, userConds, h.Offset, h.Limit)
	if err != nil {
		return nil, total, err
	}
	handler.taskUsers = userInfos

	handler.formalizeUserTask()
	return handler.userTaskInfos, total, nil
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
		return err
	}

	for _, event := range events {
		h.events[event.EntID] = event
	}
	return nil
}

func (h *Handler) GetTaskConfig(ctx context.Context, info *taskconfigmwpb.TaskConfig) (*npool.TaskConfig, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}

	if info == nil {
		taskInfo, err := taskconfigmwcli.GetTaskConfig(ctx, *h.EntID)
		if err != nil {
			return nil, err
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
		return nil, err
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
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}
	handler.taskConfigs = infos

	if err := handler.getEvents(ctx); err != nil {
		return nil, 0, err
	}

	handler.formalize()

	return handler.taskInfos, total, nil
}
