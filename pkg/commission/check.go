package commission

import (
	"context"
	"fmt"

	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
)

type checkHandler struct {
	*Handler
}

func (h *checkHandler) checkCommission(ctx context.Context) error {
	exist, err := commissionconfigmwcli.ExistCommissionConds(ctx, &commissionconfigmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt: &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid commission")
	}
	return nil
}
