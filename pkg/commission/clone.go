package commission

import (
	"context"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
)

func CloneCommissions(ctx context.Context, appID, fromGoodID, toGoodID string) error {
	return commmwcli.CloneCommissions(ctx, appID, fromGoodID, toGoodID)
}
