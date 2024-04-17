package config

import (
	"context"
	"fmt"

	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
)

type checkHandler struct {
	*Handler
}

func (h *checkHandler) checkConfig(ctx context.Context) error {
	exist, err := appconfigmwcli.ExistAppConfigConds(ctx, &appconfigmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt: &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid appconfig")
	}
	return nil
}
