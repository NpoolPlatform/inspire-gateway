package reconciliation

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	accountingmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/accounting"
	accountingmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/accounting"

	ledgermwcli "github.com/NpoolPlatform/ledger-middleware/pkg/client/ledger/v2"
	ledgerdetailmgrpb "github.com/NpoolPlatform/message/npool/ledger/mgr/v1/ledger/detail"

	ordermwpb "github.com/NpoolPlatform/message/npool/order/mw/v1/order"
	ordercli "github.com/NpoolPlatform/order-middleware/pkg/client/order"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/appgood"

	"github.com/shopspring/decimal"

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

		goodIDs := []string{}
		for _, ord := range orders {
			goodIDs = append(goodIDs, ord.GoodID)
		}

		goods, _, err := goodmwcli.GetGoods(ctx, &goodmgrpb.Conds{
			AppID: &npool.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
			GoodIDs: &npool.StringSliceVal{
				Op:    cruder.IN,
				Value: goodIDs,
			},
		}, int32(0), int32(len(goodIDs)))
		if err != nil {
			return err
		}

		goodMap := map[string]*goodmwpb.Good{}
		for _, g := range goods {
			goodMap[g.GoodID] = g
		}

		for _, order := range orders {
			good, ok := goodMap[order.GoodID]
			if !ok {
				continue
			}

			price, err := decimal.NewFromString(good.Price)
			if err != nil {
				return err
			}

			goodValue := price.Mul(decimal.NewFromInt32(int32(order.Units))).String()

			paymentAmount, err := decimal.NewFromString(order.PaymentAmount)
			if err != nil {
				return err
			}

			payWithBalance, err := decimal.NewFromString(order.PayWithBalanceAmount)
			if err != nil {
				return err
			}

			paymentAmountS := paymentAmount.Add(payWithBalance).String()

			comms, err := accountingmwcli.Accounting(ctx, &accountingmwpb.AccountingRequest{
				AppID:                  order.AppID,
				UserID:                 order.UserID,
				GoodID:                 order.GoodID,
				OrderID:                order.ID,
				PaymentID:              order.PaymentID,
				CoinTypeID:             good.CoinTypeID,
				PaymentCoinTypeID:      order.PaymentCoinTypeID,
				PaymentCoinUSDCurrency: order.PaymentCoinUSDCurrency,
				Units:                  order.Units,
				PaymentAmount:          paymentAmountS,
				GoodValue:              goodValue,
				SettleType:             good.CommissionSettleType,
			})
			if err != nil {
				logger.Sugar().Warnw("UpdateArchivement", "OrderID", order.ID, "error", err)
				continue
			}

			if len(comms) == 0 {
				continue
			}

			logger.Sugar().Warnw("UpdateArchivement", "OrderID", order.ID, "Comms", comms)

			details := []*ledgerdetailmgrpb.DetailReq{}
			ioType := ledgerdetailmgrpb.IOType_Incoming
			ioSubType := ledgerdetailmgrpb.IOSubType_Commission

			for _, comm := range comms {
				ioExtra := fmt.Sprintf(
					`{"PaymentID":"%v","OrderID":"%v","DirectContributorID":"%v","OrderUserID":"%v"}`,
					order.PaymentID,
					order.ID,
					comm.GetDirectContributorUserID(),
					order.UserID,
				)

				details = append(details, &ledgerdetailmgrpb.DetailReq{
					AppID:      &order.AppID,
					UserID:     &comm.UserID,
					CoinTypeID: &good.CoinTypeID,
					IOType:     &ioType,
					IOSubType:  &ioSubType,
					Amount:     &comm.Amount,
					IOExtra:    &ioExtra,
				})
			}

			err = ledgermwcli.BookKeeping(ctx, details)
			if err != nil {
				return err
			}
		}

		offset += int32(len(orders))
	}

	return nil
}
