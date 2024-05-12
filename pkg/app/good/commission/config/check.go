package config

import (
	"context"
	"fmt"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/good/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/good/commission/config"
)

type checkHandler struct {
	*Handler
}

func (h *checkHandler) checkConfig(ctx context.Context) error {
	exist, err := commissionconfigmwcli.ExistCommissionConfigConds(ctx, &commissionconfigmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt: &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid appgoodcommissionconfig")
	}
	return nil
}

func (h *checkHandler) checkGood(ctx context.Context) error {
	if h.AppGoodID == nil {
		return nil
	}

	exist, err := appgoodmwcli.ExistGoodConds(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid appgood")
	}
	return nil
}
