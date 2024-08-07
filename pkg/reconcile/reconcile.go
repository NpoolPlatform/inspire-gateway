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

	"github.com/shopspring/decimal"
)

type reconcileHandler struct {
	*Handler
	orderIDs  []string
	feeorders map[string][]*feeordermwpb.FeeOrder
}

func (h *reconcileHandler) reconcilePowerRentalOrder(ctx context.Context, powerRentalOrder *powerrentalordermwpb.PowerRentalOrder) error {
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

	statementReqs, err := calculatemwcli.ReconcileCalculate(ctx, &calculatemwpb.ReconcileCalculateRequest{
		OrderID: powerRentalOrder.OrderID,
	})
	if err != nil {
		return err
	}

	if len(statementReqs) == 0 {
		return nil
	}

	if err := orderstatementmwcli.UpdateStatements(ctx, statementReqs); err != nil {
		return err
	}

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
				`{"PaymentID":"%v","OrderID":"%v","OrderUserID":"%v","InspireAppConfigID":"%v","CommissionConfigID":"%v","CommissionConfigType":"%v","PaymentStatementID":"%v"}`,
				powerRentalOrder.PaymentID,
				powerRentalOrder.OrderID,
				powerRentalOrder.UserID,
				*statement.AppConfigID,
				*statement.CommissionConfigID,
				*statement.CommissionConfigType,
				*paymentStatement.EntID,
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
	powerRentalOrders := []*powerrentalordermwpb.PowerRentalOrder{}
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
			powerRentalOrders = append(powerRentalOrders, order)
			h.orderIDs = append(h.orderIDs, order.OrderID)
		}

		offset += limit
	}

	offset = 0
	for {
		childs, _, err := feeordermwcli.GetFeeOrders(ctx, &feeordermwpb.Conds{
			ParentOrderIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.orderIDs},
			PaymentType:    &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(ordertypes.PaymentType_PayWithParentOrder)},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(childs) == 0 {
			break
		}
		for _, child := range childs {
			feeorders, ok := h.feeorders[child.ParentOrderID]
			if !ok {
				feeorders = []*feeordermwpb.FeeOrder{}
			}
			feeorders = append(feeorders, child)
			h.feeorders[child.ParentOrderID] = feeorders
		}
		offset += limit
	}

	for _, order := range powerRentalOrders {
		if err := h.reconcilePowerRentalOrder(ctx, order); err != nil {
			logger.Sugar().Errorw(
				"reconcileOrders",
				"AppID", *h.AppID,
				"AppGoodID", *h.AppGoodID,
				"OrderID", order.OrderID,
				"Err", err,
			)
		}
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
		Handler:   h,
		orderIDs:  []string{},
		feeorders: map[string][]*feeordermwpb.FeeOrder{},
	}
	if err := handler.checkAppCommissionType(ctx); err != nil {
		return err
	}
	if err := handler.checkAppGoodType(ctx); err != nil {
		return err
	}
	return handler.reconcilePowerRentalOrders(ctx)
}
