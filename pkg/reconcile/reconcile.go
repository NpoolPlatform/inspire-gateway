package reconcile

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	apppowerrentalmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/powerrental"
	orderstatementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement/statement/order"
	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	calculatemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/calculate"
	ledgerstatementmwcli "github.com/NpoolPlatform/ledger-middleware/pkg/client/ledger/statement"
	goodtypes "github.com/NpoolPlatform/message/npool/basetypes/good/v1"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	ledgertypes "github.com/NpoolPlatform/message/npool/basetypes/ledger/v1"
	ordertypes "github.com/NpoolPlatform/message/npool/basetypes/order/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	apppowerrentalmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/powerrental"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
	calculatemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/calculate"
	ledgerstatementmwpb "github.com/NpoolPlatform/message/npool/ledger/mw/v2/ledger/statement"
	feeordermwpb "github.com/NpoolPlatform/message/npool/order/mw/v1/fee"
	powerrentalordermwpb "github.com/NpoolPlatform/message/npool/order/mw/v1/powerrental"
	feeordermwcli "github.com/NpoolPlatform/order-middleware/pkg/client/fee"
	powerrentalordermwcli "github.com/NpoolPlatform/order-middleware/pkg/client/powerrental"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type reconcileHandler struct {
	*Handler
}

func (h reconcileHandler) powerRentalOrderGoodValue(ctx context.Context, powerRentalOrder *powerrentalordermwpb.PowerRentalOrder) (decimal.Decimal, error) {
	offset := int32(0)
	limit := constant.DefaultRowLimit

	goodValueUSD, err := decimal.NewFromString(powerRentalOrder.GoodValueUSD)
	if err != nil {
		return decimal.NewFromInt(0), err
	}

	for {
		childs, _, err := feeordermwcli.GetFeeOrders(ctx, &feeordermwpb.Conds{
			ParentOrderID: &basetypes.StringVal{Op: cruder.EQ, Value: powerRentalOrder.EntID},
			PaymentType:   &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(ordertypes.PaymentType_PayWithParentOrder)},
		}, offset, limit)
		if err != nil {
			return decimal.NewFromInt(0), err
		}
		if len(childs) == 0 {
			break
		}
		for _, child := range childs {
			amountUSD, err := decimal.NewFromString(child.GoodValueUSD)
			if err != nil {
				return decimal.NewFromInt(0), err
			}
			goodValueUSD = goodValueUSD.Add(amountUSD)
		}
		offset += limit
	}
	return goodValueUSD, nil
}

func (h *reconcileHandler) reconcilePowerRentalOrder(ctx context.Context, powerRentalOrder *powerrentalordermwpb.PowerRentalOrder) error { //nolint
	appPowerRental, err := apppowerrentalmwcli.GetPowerRentalOnly(ctx, &apppowerrentalmwpb.Conds{
		AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: powerRentalOrder.AppID},
		AppGoodID: &basetypes.StringVal{Op: cruder.EQ, Value: powerRentalOrder.AppGoodID},
	})
	if err != nil {
		return err
	}
	if appPowerRental == nil {
		return fmt.Errorf("invalid apppowerrental")
	}

	goodValueUSD, err := h.powerRentalOrderGoodValue(ctx, powerRentalOrder)
	if err != nil {
		return err
	}

	statementReqs, err := calculatemwcli.Calculate(ctx, &calculatemwpb.CalculateRequest{
		AppID:     powerRentalOrder.AppID,
		UserID:    powerRentalOrder.UserID,
		GoodID:    powerRentalOrder.GoodID,
		AppGoodID: powerRentalOrder.AppGoodID,
		OrderID:   powerRentalOrder.OrderID,
		GoodCoinTypeID: func() string {
			uid := uuid.Nil.String()
			for _, goodCoin := range appPowerRental.GoodCoins {
				if goodCoin.Main {
					uid = goodCoin.CoinTypeID
					break
				}
			}
			return uid
		}(),
		Units:            powerRentalOrder.Units,
		PaymentAmountUSD: powerRentalOrder.PaymentAmountUSD,
		GoodValueUSD:     goodValueUSD.String(),
		SettleType:       types.SettleType_GoodOrderPayment,
		HasCommission:    powerRentalOrder.OrderType == ordertypes.OrderType_Normal,
		OrderCreatedAt:   powerRentalOrder.CreatedAt,
		Payments: func() (payments []*calculatemwpb.Payment) {
			for _, payment := range powerRentalOrder.PaymentBalances {
				payments = append(payments, &calculatemwpb.Payment{
					CoinTypeID: payment.CoinTypeID,
					Amount:     payment.Amount,
				})
			}
			for _, payment := range powerRentalOrder.PaymentTransfers {
				payments = append(payments, &calculatemwpb.Payment{
					CoinTypeID: payment.CoinTypeID,
					Amount:     payment.Amount,
				})
			}
			return
		}(),
	})
	if err != nil {
		return err
	}

	if len(statementReqs) == 0 {
		return nil
	}

	if err := orderstatementmwcli.CreateStatements(ctx, statementReqs); err != nil {
		return err
	}

	// TODO: get payment statements

	ledgerStatementReqs := []*ledgerstatementmwpb.StatementReq{}
	ioType := ledgertypes.IOType_Incoming
	ioSubType := ledgertypes.IOSubType_Commission

	for _, statement := range statementReqs {
		for _, paymentStatement := range statement.PaymentStatements {
			commission, err := decimal.NewFromString(*paymentStatement.CommissionAmount)
			if err != nil {
				return err
			}
			if commission.Cmp(decimal.NewFromInt(0)) <= 0 {
				continue
			}
			ioExtra := fmt.Sprintf(
				`{"PaymentID":"%v","OrderID":"%v","OrderUserID":"%v","InspireAppConfigID":"%v","CommissionConfigID":"%v","CommissionConfigType":"%v", "PaymentStatementID":"%v"}`,
				powerRentalOrder.PaymentID,
				powerRentalOrder.OrderID,
				powerRentalOrder.UserID,
				*statement.AppConfigID,
				*statement.CommissionConfigID,
				*statement.CommissionConfigType,
				*paymentStatement.StatementID,
			)

			ledgerStatementReqs = append(ledgerStatementReqs, &ledgerstatementmwpb.StatementReq{
				AppID:      &powerRentalOrder.AppID,
				UserID:     statement.UserID,
				CoinTypeID: paymentStatement.PaymentCoinTypeID,
				IOType:     &ioType,
				IOSubType:  &ioSubType,
				Amount:     paymentStatement.CommissionAmount,
				IOExtra:    &ioExtra,
			})
		}
	}

	if len(ledgerStatementReqs) == 0 {
		return nil
	}

	if _, err = ledgerstatementmwcli.CreateStatements(ctx, ledgerStatementReqs); err != nil {
		return err
	}

	return nil
}

func (h *reconcileHandler) reconcilePowerRentalOrders(ctx context.Context) error {
	offset := int32(0)
	limit := constant.DefaultRowLimit
	for {
		orders, _, err := powerrentalordermwcli.GetPowerRentalOrders(ctx, &powerrentalordermwpb.Conds{
			AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			AppGoodID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
			OrderTypes: &basetypes.Uint32SliceVal{Op: cruder.IN, Value: []uint32{
				uint32(ordertypes.OrderType_Normal),
				uint32(ordertypes.OrderType_Offline),
			}},
			OrderStates: &basetypes.Uint32SliceVal{Op: cruder.IN, Value: []uint32{
				uint32(ordertypes.OrderState_OrderStatePaid),
				uint32(ordertypes.OrderState_OrderStateInService),
				uint32(ordertypes.OrderState_OrderStateExpired),
			}},
			Simulate:  &basetypes.BoolVal{Op: cruder.EQ, Value: false},
			CreatedAt: &basetypes.Uint32Val{Op: cruder.GT, Value: 1714363200},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			break
		}

		for _, order := range orders {
			if err := h.reconcilePowerRentalOrder(ctx, order); err != nil {
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

func (h *reconcileHandler) checkAppGoodType(ctx context.Context) error {
	exist, err := appgoodmwcli.ExistGoodConds(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
		GoodTypes: &basetypes.Uint32SliceVal{Op: cruder.IN, Value: []uint32{
			uint32(goodtypes.GoodType_PowerRental),
			uint32(goodtypes.GoodType_LegacyPowerRental),
		}},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid powerrental")
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
	if err := handler.checkAppGoodType(ctx); err != nil {
		return err
	}
	return handler.reconcilePowerRentalOrders(ctx)
}
