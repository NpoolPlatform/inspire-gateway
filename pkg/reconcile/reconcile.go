package reconcile

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	statementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement/statement"
	calculatemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/calculate"
	ledgermwcli "github.com/NpoolPlatform/ledger-middleware/pkg/client/ledger/statement"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	ledgertypes "github.com/NpoolPlatform/message/npool/basetypes/ledger/v1"
	ordertypes "github.com/NpoolPlatform/message/npool/basetypes/order/v1"
	commonpb "github.com/NpoolPlatform/message/npool/basetypes/v1"
	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	statementmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement/statement"
	calculatemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/calculate"
	ledgerdetailmgrpb "github.com/NpoolPlatform/message/npool/ledger/mw/v2/ledger/statement"
	ordermwpb "github.com/NpoolPlatform/message/npool/order/mw/v1/order"
	ordermwcli "github.com/NpoolPlatform/order-middleware/pkg/client/order"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type reconcileHandler struct {
	*Handler
}

func (h *reconcileHandler) reconcileOrder(ctx context.Context, order *ordermwpb.Order) error { //nolint
	good, err := goodmwcli.GetGoodOnly(ctx, &goodmgrpb.Conds{
		AppID:  &commonpb.StringVal{Op: cruder.EQ, Value: order.AppID},
		GoodID: &commonpb.StringVal{Op: cruder.EQ, Value: order.GoodID},
	})
	if err != nil {
		return err
	}

	paymentAmount, err := decimal.NewFromString(order.PaymentAmount)
	if err != nil {
		return err
	}
	payWithBalance, err := decimal.NewFromString(order.PaymentBalanceAmount)
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
	currency, err := decimal.NewFromString(order.PaymentCoinUSDCurrency)
	if err != nil {
		return err
	}

	goodValue := price.Mul(untis).Div(currency).String()
	paymentAmountS := paymentAmount.Add(payWithBalance).String()

	logger.Sugar().Infow(
		"reconcileOrder",
		"AppID", order.AppID,
		"UserID", order.UserID,
		"OrderID", order.ID,
		"PaymentAmount", paymentAmountS,
		"GoodValue", goodValue,
		"CoinTypeID", good.CoinTypeID,
		"PaymentCoinTypeID", order.PaymentCoinTypeID,
		"USDCurrency", order.PaymentCoinUSDCurrency,
	)

	statements, err := calculatemwcli.Calculate(ctx, &calculatemwpb.CalculateRequest{
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
		SettleType:             types.SettleType_GoodOrderPayment,
		HasCommission:          order.OrderType == ordertypes.OrderType_Normal,
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
			"CoinTypeID", good.CoinTypeID,
			"PaymentCoinTypeID", order.PaymentCoinTypeID,
			"Error", err,
		)
		return err
	}

	if len(statements) == 0 {
		return nil
	}

	statementReqs := []*statementmwpb.StatementReq{}
	for _, statement := range statements {
		req := &statementmwpb.StatementReq{
			AppID:                  &statement.AppID,
			UserID:                 &statement.UserID,
			GoodID:                 &statement.GoodID,
			OrderID:                &statement.OrderID,
			SelfOrder:              &statement.SelfOrder,
			PaymentID:              &statement.PaymentID,
			CoinTypeID:             &statement.CoinTypeID,
			PaymentCoinTypeID:      &statement.PaymentCoinTypeID,
			PaymentCoinUSDCurrency: &statement.PaymentCoinUSDCurrency,
			Units:                  &statement.Units,
			Amount:                 &statement.Amount,
			USDAmount:              &statement.USDAmount,
			Commission:             &statement.Commission,
		}
		if _, err := uuid.Parse(statement.DirectContributorID); err == nil {
			req.DirectContributorID = &statement.DirectContributorID
		}
		statementReqs = append(statementReqs, req)
	}

	_, err = statementmwcli.CreateStatements(ctx, statementReqs)
	if err != nil {
		logger.Sugar().Infow(
			"reconcileOrder",
			"AppID", order.AppID,
			"UserID", order.UserID,
			"OrderID", order.ID,
			"PaymentAmount", paymentAmountS,
			"GoodValue", goodValue,
			"CoinTypeID", good.CoinTypeID,
			"PaymentCoinTypeID", order.PaymentCoinTypeID,
			"Error", err,
		)
		return err
	}

	details := []*ledgerdetailmgrpb.StatementReq{}
	ioType := ledgertypes.IOType_Incoming
	ioSubType := ledgertypes.IOSubType_Commission

	logger.Sugar().Infow(
		"reconcileOrder",
		"AppID", order.AppID,
		"UserID", order.UserID,
		"OrderID", order.ID,
		"PaymentAmount", paymentAmountS,
		"GoodValue", goodValue,
		"CoinTypeID", good.CoinTypeID,
		"PaymentCoinTypeID", order.PaymentCoinTypeID,
	)

	for _, statement := range statements {
		commission, err := decimal.NewFromString(statement.Commission)
		if err != nil {
			return err
		}
		if commission.Cmp(decimal.NewFromInt(0)) <= 0 {
			continue
		}

		logger.Sugar().Infow(
			"reconcileOrder",
			"AppID", statement.AppID,
			"UserID", statement.UserID,
			"Amount", statement.Amount,
			"DirectContributorUserID", statement.DirectContributorID,
			"OrderID", order.ID,
			"OrderUserID", order.UserID,
		)

		ioExtra := fmt.Sprintf(
			`{"PaymentID":"%v","OrderID":"%v","DirectContributorID":"%v","OrderUserID":"%v"}`,
			order.PaymentID,
			order.ID,
			statement.GetDirectContributorID(),
			order.UserID,
		)

		details = append(details, &ledgerdetailmgrpb.StatementReq{
			AppID:      &order.AppID,
			UserID:     &statement.UserID,
			CoinTypeID: &order.PaymentCoinTypeID,
			IOType:     &ioType,
			IOSubType:  &ioSubType,
			Amount:     &statement.Commission,
			IOExtra:    &ioExtra,
		})
	}

	if len(details) == 0 {
		return nil
	}

	if _, err = ledgermwcli.CreateStatements(ctx, details); err != nil {
		logger.Sugar().Infow(
			"reconcileOrder",
			"AppID", order.AppID,
			"UserID", order.UserID,
			"GoodID", order.GoodID,
			"OrderID", order.ID,
			"PaymentAmount", paymentAmountS,
			"GoodValue", goodValue,
			"CoinTypeID", good.CoinTypeID,
			"PaymentCoinTypeID", order.PaymentCoinTypeID,
			"Error", err,
		)
		return err
	}

	return nil
}

func (h *reconcileHandler) reconcileOrders(ctx context.Context, orderType ordertypes.OrderType) error {
	offset := int32(0)
	limit := constant.DefaultRowLimit
	for {
		orders, _, err := ordermwcli.GetOrders(
			ctx,
			&ordermwpb.Conds{
				AppID:     &commonpb.StringVal{Op: cruder.EQ, Value: *h.AppID},
				GoodID:    &commonpb.StringVal{Op: cruder.EQ, Value: *h.GoodID},
				OrderType: &commonpb.Uint32Val{Op: cruder.EQ, Value: uint32(orderType)},
				OrderStates: &commonpb.Uint32SliceVal{
					Op: cruder.IN,
					Value: []uint32{
						uint32(ordertypes.OrderState_OrderStatePaid),
						uint32(ordertypes.OrderState_OrderStateInService),
						uint32(ordertypes.OrderState_OrderStateExpired),
					},
				},
			},
			offset,
			limit,
		)
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			break
		}

		for _, order := range orders {
			if err := h.reconcileOrder(ctx, order); err != nil {
				logger.Sugar().Errorw(
					"reconcileOrders",
					"AppID", *h.AppID,
					"GoodID", *h.GoodID,
					"OrderID", order.ID,
					"Err", err,
				)
			}
		}

		offset += limit
	}
	return nil
}

func (h *Handler) Reconcile(ctx context.Context) error {
	if h.AppID == nil {
		return fmt.Errorf("invalid appid")
	}
	if h.GoodID == nil {
		return fmt.Errorf("invalid goodid")
	}
	handler := &reconcileHandler{
		Handler: h,
	}
	if err := handler.reconcileOrders(ctx, ordertypes.OrderType_Normal); err != nil {
		return err
	}
	return handler.reconcileOrders(ctx, ordertypes.OrderType_Offline)
}
