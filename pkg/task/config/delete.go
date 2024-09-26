package config

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	configmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"
	configmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/config"
)

func (h *Handler) DeleteTaskConfig(ctx context.Context) (*npool.TaskConfig, error) {
	info, err := configmwcli.GetTaskConfigOnly(ctx, &configmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, wlog.WrapError(err)
	}
	if info == nil {
		return nil, wlog.Errorf("invalid taskconfig")
	}
	h.AppID = &info.AppID
	exist, err := configmwcli.ExistTaskConfigConds(ctx, &configmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		LastTaskID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, wlog.WrapError(err)
	}
	if exist {
		return nil, wlog.Errorf("invalid lasttaskid")
	}

	if err := configmwcli.DeleteTaskConfig(ctx, h.ID, h.EntID); err != nil {
		return nil, wlog.WrapError(err)
	}

	return h.GetTaskConfig(ctx, info)
}
