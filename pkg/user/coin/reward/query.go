package reward

import (
	"context"
	"fmt"

	usercoinrewardmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/user/coin/reward"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	usercoinrewardmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/user/coin/reward"
)

func (h *Handler) GetUserCoinReward(ctx context.Context) (*usercoinrewardmwpb.UserCoinReward, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}

	info, err := usercoinrewardmwcli.GetUserCoinReward(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	return info, nil
}

func (h *Handler) GetUserCoinRewards(ctx context.Context) ([]*usercoinrewardmwpb.UserCoinReward, uint32, error) {
	conds := &usercoinrewardmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}

	return usercoinrewardmwcli.GetUserCoinRewards(ctx, conds, h.Offset, h.Limit)
}
