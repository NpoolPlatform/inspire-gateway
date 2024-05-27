package config

import (
	"context"
	"fmt"

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
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid taskconfig")
	}

	if err := configmwcli.DeleteTaskConfig(ctx, h.ID, h.EntID); err != nil {
		return nil, err
	}

	return h.GetTaskConfig(ctx, info)
}
