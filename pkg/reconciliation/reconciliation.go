package reconciliation

import (
	"context"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	accountingmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/accounting"
	accountingmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/accounting"

	commmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"

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

		// TODO: get good list

		for _, order := range orders {
			comms, err := accountingmwcli.Accounting(ctx, &accountingmwpb.AccountingRequest{
				AppID:     order.AppID,
				UserID:    order.UserID,
				GoodID:    order.GoodID,
				OrderID:   order.ID,
				PaymentID: order.PaymentID,
				// CoinTypeID:             order.CoinTypeID,
				PaymentCoinTypeID:      order.PaymentCoinTypeID,
				PaymentCoinUSDCurrency: order.PaymentCoinUSDCurrency,
				Units:                  order.Units,
				PaymentAmount:          order.PaymentAmount, // PayWithBalanceAmount
				//GoodValue:              order.GoodValue,
				SettleType: commmgrpb.SettleType_GoodOrderPercent,
			})
			if err != nil {
				logger.Sugar().Warnw("UpdateArchivement", "OrderID", order.ID, "error", err)
				continue
			}

			logger.Sugar().Warnw("UpdateArchivement", "OrderID", order.ID, "Comms", comms)
			// TODO: update comm to ledger
		}

		offset += int32(len(orders))
	}

	return nil
}
