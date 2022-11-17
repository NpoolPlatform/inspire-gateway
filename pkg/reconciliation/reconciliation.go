package reconciliation

import (
	"context"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	archivement "github.com/NpoolPlatform/staker-manager/pkg/archivement"
	commission "github.com/NpoolPlatform/staker-manager/pkg/commission"

	ordermwpb "github.com/NpoolPlatform/message/npool/order/mw/v1/order"
	ordercli "github.com/NpoolPlatform/order-middleware/pkg/client/order"

	"github.com/NpoolPlatform/message/npool"
)

// nolint
func UpdateArchivement(ctx context.Context, appID, userID string) error {
	offset := int32(0)
	limit := int32(1000) //nolint // Mock variable now

	for {
		orders, _, err := ordercli.GetOrders(ctx, &ordermwpb.Conds{
			AppID: &npool.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
			UserID: &npool.StringVal{
				Op:    cruder.EQ,
				Value: userID,
			},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			return nil
		}

		for _, order := range orders {
			if err := commission.CalculateCommission(ctx, order.ID, false); err != nil {
				logger.Sugar().Warnw("UpdateArchivement", "OrderID", order.ID, "error", err)
				continue
			}
			if err := archivement.CalculateArchivement(ctx, order.ID); err != nil {
				logger.Sugar().Warnw("UpdateArchivement", "OrderID", order.ID, "error", err)
			}
		}

		offset += int32(len(orders))
		// Only mock, so just return
		return nil
	}
}
