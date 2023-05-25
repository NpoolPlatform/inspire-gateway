package commission

import (
	"context"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	commmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"
)

func CloneCommissions(ctx context.Context, appID, fromGoodID, toGoodID, value string, settleType commmgrpb.SettleType) error {
	return commmwcli.CloneCommissions(ctx, appID, fromGoodID, toGoodID, value, settleType)
}
