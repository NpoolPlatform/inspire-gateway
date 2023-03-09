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

	ordermgrpb "github.com/NpoolPlatform/message/npool/order/mgr/v1/order"

	"github.com/shopspring/decimal"

	"github.com/NpoolPlatform/message/npool"
)

func reconcileOrder(ctx context.Context, order *ordermwpb.Order) error { //nolint
	good, err := goodmwcli.GetGoodOnly(ctx, &goodmgrpb.Conds{
		AppID:  &npool.StringVal{Op: cruder.EQ, Value: order.AppID},
		GoodID: &npool.StringVal{Op: cruder.EQ, Value: order.GoodID},
	})
	if err != nil {
		return err
	}

	paymentAmount, err := decimal.NewFromString(order.PaymentAmount)
	if err != nil {
		return err
	}
	payWithBalance, err := decimal.NewFromString(order.PayWithBalanceAmount)
	if err != nil {
		return err
	}

	price, err := decimal.NewFromString(good.Price)
	if err != nil {
		return err
	}
	untis, err := decimal.NewFromString(order.Units)
	if err != nil {
		return err
	}

	goodValue := price.Mul(untis).String()
	paymentAmountS := paymentAmount.Add(payWithBalance).String()

	logger.Sugar().Infow(
		"reconcileOrder",
		"AppID", order.AppID,
		"UserID", order.UserID,
		"OrderID", order.ID,
		"PaymentAmount", paymentAmountS,
		"GoodValue", goodValue,
		"SettleType", good.CommissionSettleType,
		"CoinTypeID", good.CoinTypeID,
		"PaymentCoinTypeID", order.PaymentCoinTypeID,
		"USDCurrency", order.PaymentCoinUSDCurrency,
	)

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
		HasCommission:          order.OrderType == ordermgrpb.OrderType_Normal,
		OrderCreatedAt:         order.CreatedAt,
	})
	if err != nil {
		logger.Sugar().Infow(
			"reconcileOrder",
			"AppID", order.AppID,
			"UserID", order.UserID,
			"OrderID", order.ID,
			"PaymentAmount", paymentAmountS,
			"GoodValue", goodValue,
			"SettleType", good.CommissionSettleType,
			"CoinTypeID", good.CoinTypeID,
			"PaymentCoinTypeID", order.PaymentCoinTypeID,
			"Error", err,
		)
		return err
	}

	if len(comms) == 0 {
		return nil
	}

	details := []*ledgerdetailmgrpb.DetailReq{}
	ioType := ledgerdetailmgrpb.IOType_Incoming
	ioSubType := ledgerdetailmgrpb.IOSubType_Commission

	logger.Sugar().Infow(
		"reconcileOrder",
		"AppID", order.AppID,
		"UserID", order.UserID,
		"OrderID", order.ID,
		"PaymentAmount", paymentAmountS,
		"GoodValue", goodValue,
		"SettleType", good.CommissionSettleType,
		"CoinTypeID", good.CoinTypeID,
		"PaymentCoinTypeID", order.PaymentCoinTypeID,
	)

	for _, comm := range comms {
		logger.Sugar().Infow(
			"reconcileOrder",
			"AppID", comm.AppID,
			"UserID", comm.UserID,
			"Amount", comm.Amount,
			"DirectContributorUserID", comm.DirectContributorUserID,
			"OrderID", order.ID,
			"OrderUserID", order.UserID,
		)

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
			CoinTypeID: &order.PaymentCoinTypeID,
			IOType:     &ioType,
			IOSubType:  &ioSubType,
			Amount:     &comm.Amount,
			IOExtra:    &ioExtra,
		})
	}

	err = ledgermwcli.BookKeeping(ctx, details)
	if err != nil {
		logger.Sugar().Infow(
			"reconcileOrder",
			"AppID", order.AppID,
			"UserID", order.UserID,
			"OrderID", order.ID,
			"PaymentAmount", paymentAmountS,
			"GoodValue", goodValue,
			"SettleType", good.CommissionSettleType,
			"CoinTypeID", good.CoinTypeID,
			"PaymentCoinTypeID", order.PaymentCoinTypeID,
			"Error", err,
		)
		return err
	}

	return nil
}

func reconcileOrders(ctx context.Context, conds *ordermwpb.Conds, offset, limit int32) (bool, error) {
	orders, _, err := ordercli.GetOrders(ctx, conds, offset, limit)
	if err != nil {
		return false, err
	}

	logger.Sugar().Infow(
		"reconcileOrders",
		"Orders", len(orders),
	)

	if len(orders) == 0 {
		return true, nil
	}

	for _, order := range orders {
		if err := reconcileOrder(ctx, order); err != nil {
			logger.Sugar().Errorw(
				"reconcileOrders",
				"AppID", order.AppID,
				"UserID", order.UserID,
				"OrderID", order.ID,
				"Error", err,
			)
			return true, err
		}
	}

	return false, nil
}

func reconcileTypedOrders(ctx context.Context, appID, userID string, orderType ordermgrpb.OrderType) error {
	logger.Sugar().Infow(
		"reconcileTypedOrders",
		"AppID", appID,
		"UserID", userID,
		"OrderType", orderType,
	)

	offset := int32(0)
	const limit = int32(100)

	for {
		finish, err := reconcileOrders(ctx, &ordermwpb.Conds{
			AppID:  &npool.StringVal{Op: cruder.EQ, Value: appID},
			UserID: &npool.StringVal{Op: cruder.EQ, Value: userID},
			Type:   &npool.Uint32Val{Op: cruder.EQ, Value: uint32(orderType)},
			States: &npool.Uint32SliceVal{
				Op: cruder.IN,
				Value: []uint32{
					uint32(ordermgrpb.OrderState_Paid),
					uint32(ordermgrpb.OrderState_InService),
					uint32(ordermgrpb.OrderState_Expired),
				},
			},
		}, offset, limit)
		if err != nil {
			logger.Sugar().Errorw(
				"reconcileTypeOrders",
				"AppID", appID,
				"UserID", userID,
				"Type", orderType,
				"Error", err,
			)
			return err
		}
		if finish {
			break
		}

		offset += limit
	}

	return nil
}

func Reconcile(ctx context.Context, appID, userID string) error {
	if err := reconcileTypedOrders(ctx, appID, userID, ordermgrpb.OrderType_Normal); err != nil {
		return err
	}
	return reconcileTypedOrders(ctx, appID, userID, ordermgrpb.OrderType_Offline)
}
