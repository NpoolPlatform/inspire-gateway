package reward

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	userrewardmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/user/reward"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	userrewardmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/user/reward"
)

func (h *Handler) GetUserReward(ctx context.Context) (*userrewardmwpb.UserReward, error) {
	if h.EntID == nil {
		return nil, wlog.Errorf("invalid entid")
	}

	info, err := userrewardmwcli.GetUserReward(ctx, *h.EntID)
	if err != nil {
		return nil, wlog.WrapError(err)
	}
	if info == nil {
		return nil, nil
	}

	return info, nil
}

func (h *Handler) GetUserRewards(ctx context.Context) ([]*userrewardmwpb.UserReward, uint32, error) {
	conds := &userrewardmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}

	return userrewardmwcli.GetUserRewards(ctx, conds, h.Offset, h.Limit)
}
