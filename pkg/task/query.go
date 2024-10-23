package config

import (
	"context"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	taskconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/config"
	taskusermwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/user"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	taskconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/config"
	taskusermwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/user"
)

type queryHandler struct {
	*Handler
	taskConfigs   []*taskconfigmwpb.TaskConfig
	taskUsers     []*taskusermwpb.TaskUser
	events        map[string]*eventmwpb.Event
	userTaskInfos []*npool.UserTask
}

func (h *queryHandler) getEvents(ctx context.Context) error {
	infos, _, err := eventmwcli.GetEvents(ctx, &eventmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return wlog.WrapError(err)
	}
	for _, event := range infos {
		h.events[event.EntID] = event
	}
	return nil
}

func (h *queryHandler) calculateNextStartAt(taskConfig *taskconfigmwpb.TaskConfig, taskUsers []*taskusermwpb.TaskUser) uint32 {
	if taskConfig.MaxRewardCount == uint32(len(taskUsers)) {
		return 0
	}
	finishedAt := uint32(0)
	for _, taskUser := range taskUsers {
		if taskUser.CreatedAt > finishedAt {
			finishedAt = taskUser.CreatedAt
		}
	}
	if !taskConfig.IntervalReset {
		return finishedAt + taskConfig.CooldownSecond
	}
	now := uint32(time.Now().Unix())
	intervalTime := int32(now / taskConfig.IntervalResetSecond)
	startTime := uint32(intervalTime) * taskConfig.IntervalResetSecond
	intervalTaskCount := uint32(0)
	for _, taskUser := range taskUsers {
		if taskUser.CreatedAt > startTime {
			intervalTaskCount++
		}
	}
	if intervalTaskCount == 0 {
		return startTime
	}
	if intervalTaskCount == taskConfig.MaxIntervalRewardCount {
		return startTime + taskConfig.IntervalResetSecond
	}
	if intervalTaskCount < taskConfig.MaxIntervalRewardCount {
		return finishedAt + taskConfig.CooldownSecond
	}

	return 0
}

func (h *queryHandler) formalizeUserTask() {
	taskCountMap := map[string]uint32{}
	taskMap := map[string][]*taskusermwpb.TaskUser{}
	for _, taskUser := range h.taskUsers {
		taskUserCount, ok := taskCountMap[taskUser.TaskID]
		if ok {
			taskUserCount++
			taskCountMap[taskUser.TaskID] = taskUserCount
			continue
		}
		taskCountMap[taskUser.TaskID] = uint32(1)

		taskUsers, ok := taskMap[taskUser.TaskID]
		if ok {
			taskUsers := append(taskUsers, taskUser)
			taskMap[taskUser.TaskID] = taskUsers
			continue
		}
		taskMap[taskUser.TaskID] = []*taskusermwpb.TaskUser{taskUser}
	}
	for _, comm := range h.taskConfigs {
		_, ok := h.events[comm.EventID]
		if !ok {
			continue
		}
		taskState := types.TaskState_NotStarted
		rewardState := types.RewardState_UnIssued
		taskUserCount, ok := taskCountMap[comm.EntID]
		if ok {
			taskState = types.TaskState_Done
			rewardState = types.RewardState_Issued
		}
		nextStartAt := uint32(0)
		taskUsers, ok := taskMap[comm.EntID]
		if ok {
			nextStartAt = h.calculateNextStartAt(comm, taskUsers)
		}
		h.userTaskInfos = append(h.userTaskInfos, &npool.UserTask{
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
			CooldownSecond:   comm.CooldownSecond,
			CompletionTimes:  taskUserCount,
			TaskState:        taskState,
			RewardState:      rewardState,
			NextStartAt:      nextStartAt,
			CreatedAt:        comm.CreatedAt,
			UpdatedAt:        comm.UpdatedAt,
		})
	}
}

func (h *Handler) GetUserTask(ctx context.Context) (*npool.UserTask, error) {
	if h.EntID == nil {
		return nil, wlog.Errorf("invalid entid")
	}

	info, err := taskconfigmwcli.GetTaskConfig(ctx, *h.EntID)
	if err != nil {
		return nil, wlog.WrapError(err)
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:       h,
		taskConfigs:   []*taskconfigmwpb.TaskConfig{info},
		taskUsers:     []*taskusermwpb.TaskUser{},
		userTaskInfos: []*npool.UserTask{},
		events:        map[string]*eventmwpb.Event{},
	}
	if err := handler.getEvents(ctx); err != nil {
		return nil, err
	}

	handler.formalizeUserTask()
	if len(handler.userTaskInfos) == 0 {
		return nil, nil
	}

	return handler.userTaskInfos[0], nil
}

func (h *Handler) GetUserTasks(ctx context.Context) ([]*npool.UserTask, uint32, error) {
	handler := &queryHandler{
		Handler:       h,
		userTaskInfos: []*npool.UserTask{},
		taskConfigs:   []*taskconfigmwpb.TaskConfig{},
		taskUsers:     []*taskusermwpb.TaskUser{},
		events:        map[string]*eventmwpb.Event{},
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

	userConds := &taskusermwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
	}

	userInfos, _, err := taskusermwcli.GetTaskUsers(ctx, userConds, h.Offset, h.Limit)
	if err != nil {
		return nil, total, wlog.WrapError(err)
	}
	handler.taskUsers = userInfos
	if err := handler.getEvents(ctx); err != nil {
		return nil, 0, err
	}

	handler.formalizeUserTask()
	return handler.userTaskInfos, total, nil
}
