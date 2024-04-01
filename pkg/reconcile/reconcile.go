package reconcile

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	achievementstatementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement/statement"
	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	calculatemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/calculate"
	ledgerstatementmwcli "github.com/NpoolPlatform/ledger-middleware/pkg/client/ledger/statement"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	ledgertypes "github.com/NpoolPlatform/message/npool/basetypes/ledger/v1"
	ordertypes "github.com/NpoolPlatform/message/npool/basetypes/order/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	achievementstatementmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement/statement"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
	calculatemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/calculate"
	ledgerstatementmwpb "github.com/NpoolPlatform/message/npool/ledger/mw/v2/ledger/statement"
	ordermwpb "github.com/NpoolPlatform/message/npool/order/mw/v1/order"
	ordermwcli "github.com/NpoolPlatform/order-middleware/pkg/client/order"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type reconcileHandler struct {
	*Handler
}

//nolint:gocritic
func (h reconcileHandler) orderGoodValue(ctx context.Context, order *ordermwpb.Order) (decimal.Decimal, decimal.Decimal, error) {
	offset := int32(0)
	limit := constant.DefaultRowLimit

	goodValue, err := decimal.NewFromString(order.GoodValue)
	if err != nil {
		return decimal.NewFromInt(0), decimal.NewFromInt(0), err
	}
	goodValueUSD, err := decimal.NewFromString(order.GoodValueUSD)
	if err != nil {
		return decimal.NewFromInt(0), decimal.NewFromInt(0), err
	}

	for {
		childs, _, err := ordermwcli.GetOrders(ctx, &ordermwpb.Conds{
			ParentOrderID: &basetypes.StringVal{Op: cruder.EQ, Value: order.EntID},
			PaymentType:   &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(ordertypes.PaymentType_PayWithParentOrder)},
		}, offset, limit)
		if err != nil {
			return decimal.NewFromInt(0), decimal.NewFromInt(0), err
		}
		if len(childs) == 0 {
			break
		}
		for _, child := range childs {
			amount, err := decimal.NewFromString(child.GoodValue)
			if err != nil {
				return decimal.NewFromInt(0), decimal.NewFromInt(0), err
			}
			amountUSD, err := decimal.NewFromString(child.GoodValueUSD)
			if err != nil {
				return decimal.NewFromInt(0), decimal.NewFromInt(0), err
			}
			goodValue = goodValue.Add(amount)
			goodValueUSD = goodValueUSD.Add(amountUSD)
		}
		offset += limit
	}
	return goodValue, goodValueUSD, nil
}

func (h *reconcileHandler) reconcileOrder(ctx context.Context, order *ordermwpb.Order) error { //nolint
	good, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: order.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: order.AppGoodID},
	})
	if err != nil {
		return err
	}
	if good == nil {
		return fmt.Errorf("invalid good")
	}

	goodValue, goodValueUSD, err := h.orderGoodValue(ctx, order)
	if err != nil {
		return err
	}

	statements, err := calculatemwcli.Calculate(ctx, &calculatemwpb.CalculateRequest{
		AppID:                  order.AppID,
		UserID:                 order.UserID,
		GoodID:                 order.GoodID,
		AppGoodID:              order.AppGoodID,
		OrderID:                order.EntID,
		PaymentID:              order.PaymentID,
		CoinTypeID:             good.CoinTypeID,
		PaymentCoinTypeID:      order.PaymentCoinTypeID,
		PaymentCoinUSDCurrency: order.CoinUSDCurrency,
		Units:                  order.Units,
		PaymentAmount:          order.PaymentAmount,
		GoodValue:              goodValue.String(),
		GoodValueUSD:           goodValueUSD.String(),
		SettleType:             types.SettleType_GoodOrderPayment,
		HasCommission:          order.OrderType == ordertypes.OrderType_Normal,
		OrderCreatedAt:         order.CreatedAt,
	})
	if err != nil {
		return err
	}

	if len(statements) == 0 {
		return nil
	}

	achievementStatementReqs := []*achievementstatementmwpb.StatementReq{}
	for _, statement := range statements {
		req := &achievementstatementmwpb.StatementReq{
			AppID:                  &statement.AppID,
			UserID:                 &statement.UserID,
			GoodID:                 &statement.GoodID,
			AppGoodID:              &statement.AppGoodID,
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
			AppConfigID:            &statement.AppConfigID,
			CommissionConfigID:     &statement.CommissionConfigID,
			CommissionConfigType:   &statement.CommissionConfigType,
		}
		if _, err := uuid.Parse(statement.DirectContributorID); err == nil {
			req.DirectContributorID = &statement.DirectContributorID
		}
		achievementStatementReqs = append(achievementStatementReqs, req)
	}

	_, err = achievementstatementmwcli.CreateStatements(ctx, achievementStatementReqs)
	if err != nil {
		return err
	}

	ledgerStatementReqs := []*ledgerstatementmwpb.StatementReq{}
	ioType := ledgertypes.IOType_Incoming
	ioSubType := ledgertypes.IOSubType_Commission

	for _, statement := range statements {
		commission, err := decimal.NewFromString(statement.Commission)
		if err != nil {
			return err
		}
		if commission.Cmp(decimal.NewFromInt(0)) <= 0 {
			continue
		}
		ioExtra := fmt.Sprintf(
			`{"PaymentID":"%v","OrderID":"%v","DirectContributorID":"%v","OrderUserID":"%v"}`,
			order.PaymentID,
			order.EntID,
			statement.GetDirectContributorID(),
			order.UserID,
		)

		ledgerStatementReqs = append(ledgerStatementReqs, &ledgerstatementmwpb.StatementReq{
			AppID:      &order.AppID,
			UserID:     &statement.UserID,
			CoinTypeID: &order.PaymentCoinTypeID,
			IOType:     &ioType,
			IOSubType:  &ioSubType,
			Amount:     &statement.Commission,
			IOExtra:    &ioExtra,
		})
	}

	if len(ledgerStatementReqs) == 0 {
		return nil
	}

	if _, err = ledgerstatementmwcli.CreateStatements(ctx, ledgerStatementReqs); err != nil {
		return err
	}

	return nil
}

func (h *reconcileHandler) checkAppCommissionType(ctx context.Context) error {
	appConfig, err := appconfigmwcli.GetAppConfigOnly(ctx, &appconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(0)},
	})
	if err != nil {
		return err
	}
	if appConfig == nil {
		return fmt.Errorf("invalid inspire appconfig")
	}
	if appConfig.CommissionType != types.CommissionType_LegacyCommission {
		return fmt.Errorf("invalid commissiontype")
	}
	return nil
}

func (h *reconcileHandler) reconcileOrders(ctx context.Context, orderType ordertypes.OrderType) error {
	offset := int32(0)
	limit := constant.DefaultRowLimit
	simulate := false
	for {
		orders, _, err := ordermwcli.GetOrders(ctx, &ordermwpb.Conds{
			AppID:       &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			AppGoodID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
			OrderType:   &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(orderType)},
			PaymentType: &basetypes.Uint32Val{Op: cruder.NEQ, Value: uint32(ordertypes.PaymentType_PayWithParentOrder)},
			OrderStates: &basetypes.Uint32SliceVal{Op: cruder.IN, Value: []uint32{
				uint32(ordertypes.OrderState_OrderStatePaid),
				uint32(ordertypes.OrderState_OrderStateInService),
				uint32(ordertypes.OrderState_OrderStateExpired),
			}},
			Simulate: &basetypes.BoolVal{Op: cruder.EQ, Value: simulate},
		}, offset, limit)
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
					"AppGoodID", *h.AppGoodID,
					"OrderID", order.EntID,
					"Err", err,
				)
			}
		}

		offset += limit
	}
	return nil
}

func (h *Handler) Reconcile(ctx context.Context) error {
	handler := &reconcileHandler{
		Handler: h,
	}
	if err := handler.checkAppCommissionType(ctx); err != nil {
		return err
	}
	if err := handler.reconcileOrders(ctx, ordertypes.OrderType_Normal); err != nil {
		return err
	}
	return handler.reconcileOrders(ctx, ordertypes.OrderType_Offline)
}
