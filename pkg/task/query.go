package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	taskconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/config"
	taskusermwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/user"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task"
	taskconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/config"
	taskusermwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/user"
)

type queryHandler struct {
	*Handler
	taskConfigs   []*taskconfigmwpb.TaskConfig
	taskUsers     []*taskusermwpb.TaskUser
	userTaskInfos []*npool.UserTask
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
			CooldownSecord:   comm.CooldownSecord,
			CompletionTimes:  taskUserCount,
			TaskState:        taskState,
			RewardState:      rewardState,
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

	handler.formalizeUserTask()
	return handler.userTaskInfos, total, nil
}
