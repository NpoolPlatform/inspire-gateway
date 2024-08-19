//nolint:funlen
package migrator

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	servicename "github.com/NpoolPlatform/inspire-gateway/pkg/servicename"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db/ent"
	entgoodcoinachievement "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/goodcoinachievement"
	entorderpaymentstatement "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/orderpaymentstatement"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	ordertypes "github.com/NpoolPlatform/message/npool/basetypes/order/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	powerrentalordermwpb "github.com/NpoolPlatform/message/npool/order/mw/v1/powerrental"

	powerrentalordermwcli "github.com/NpoolPlatform/order-middleware/pkg/client/powerrental"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	keyServiceID = "serviceid"
)

func lockKey() string {
	serviceID := config.GetStringValueWithNameSpace(servicename.ServiceDomain, keyServiceID)
	return fmt.Sprintf("%v:%v", basetypes.Prefix_PrefixMigrate, serviceID)
}

//nolint:gocyclo
func migrateAchievement(ctx context.Context, tx *ent.Tx) error {
	type MyAchievement struct {
		EntID uuid.UUID `json:"ent_id"`
	}

	rows, err := tx.QueryContext(ctx, "select ent_id from good_achievements")
	if err != nil {
		return wlog.WrapError(err)
	}
	goodAchievements := map[uuid.UUID]bool{}
	for rows.Next() {
		goodAchievement := &MyAchievement{}
		if err := rows.Scan(&goodAchievement.EntID); err != nil {
			return wlog.WrapError(err)
		}
		goodAchievements[goodAchievement.EntID] = true
	}

	rows, err = tx.QueryContext(ctx, "select ent_id from good_coin_achievements")
	if err != nil {
		return wlog.WrapError(err)
	}
	goodCoinAchievements := map[uuid.UUID]bool{}
	for rows.Next() {
		goodCoinAchievement := &MyAchievement{}
		if err := rows.Scan(&goodCoinAchievement.EntID); err != nil {
			return wlog.WrapError(err)
		}
		goodCoinAchievements[goodCoinAchievement.EntID] = true
	}

	rows, err = tx.QueryContext(ctx, "select ent_id,app_id,user_id,good_id,app_good_id,coin_type_id,total_units_v1,self_units_v1,total_amount,self_amount,total_commission,self_commission,created_at,updated_at from archivement_generals where deleted_at = 0") //nolint
	if err != nil {
		return err
	}

	type Achievement struct {
		EntID           uuid.UUID       `json:"ent_id"`
		AppID           uuid.UUID       `json:"app_id"`
		UserID          uuid.UUID       `json:"user_id"`
		GoodID          uuid.UUID       `json:"good_id"`
		AppGoodID       uuid.UUID       `json:"app_good_id"`
		CoinTypeID      uuid.UUID       `json:"coin_type_id"`
		TotalUnitsV1    decimal.Decimal `json:"total_units_v1"`
		SelfUnitsV1     decimal.Decimal `json:"self_units_v1"`
		TotalAmount     decimal.Decimal `json:"total_amount"`
		SelfAmount      decimal.Decimal `json:"self_amount"`
		TotalCommission decimal.Decimal `json:"total_commission"`
		SelfCommission  decimal.Decimal `json:"self_commission"`
		CreatedAt       uint32          `json:"created_at"`
		UpdatedAt       uint32          `json:"updated_at"`
	}
	achievements := []*Achievement{}
	for rows.Next() {
		achievement := &Achievement{}
		if err := rows.Scan(
			&achievement.EntID,
			&achievement.AppID,
			&achievement.UserID,
			&achievement.GoodID,
			&achievement.AppGoodID,
			&achievement.CoinTypeID,
			&achievement.TotalUnitsV1,
			&achievement.SelfUnitsV1,
			&achievement.TotalAmount,
			&achievement.SelfAmount,
			&achievement.TotalCommission,
			&achievement.SelfCommission,
			&achievement.CreatedAt,
			&achievement.UpdatedAt,
		); err != nil {
			return err
		}
		achievements = append(achievements, achievement)
	}

	for _, achievement := range achievements {
		_, ok := goodAchievements[achievement.EntID]
		if !ok {
			if _, err := tx.
				GoodAchievement.
				Create().
				SetEntID(achievement.EntID).
				SetAppID(achievement.AppID).
				SetUserID(achievement.UserID).
				SetGoodID(achievement.GoodID).
				SetAppGoodID(achievement.AppGoodID).
				SetTotalUnits(achievement.TotalUnitsV1).
				SetSelfUnits(achievement.SelfUnitsV1).
				SetTotalAmountUsd(achievement.TotalAmount).
				SetSelfAmountUsd(achievement.SelfAmount).
				SetTotalCommissionUsd(achievement.TotalCommission).
				SetSelfCommissionUsd(achievement.SelfCommission).
				SetCreatedAt(achievement.CreatedAt).
				SetUpdatedAt(achievement.UpdatedAt).
				Save(ctx); err != nil {
				return err
			}
		}

		_, ok = goodCoinAchievements[achievement.EntID]
		if !ok {
			// need merge multi records into one exist record if cointypeid & userid is same
			goodCoinAchievement, err := tx.
				GoodCoinAchievement.
				Query().
				Where(
					entgoodcoinachievement.UserID(achievement.UserID),
					entgoodcoinachievement.GoodCoinTypeID(achievement.CoinTypeID),
				).
				Only(ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return err
				}
			}

			if goodCoinAchievement != nil { // update
				totalUnit := goodCoinAchievement.TotalUnits.Add(achievement.TotalUnitsV1)
				selfUnits := goodCoinAchievement.SelfUnits.Add(achievement.SelfUnitsV1)
				totalAmountUsd := goodCoinAchievement.TotalAmountUsd.Add(achievement.TotalAmount)
				selfAmountUsd := goodCoinAchievement.SelfAmountUsd.Add(achievement.SelfAmount)
				totalCommissionUsd := goodCoinAchievement.TotalCommissionUsd.Add(achievement.TotalCommission)
				selfCommissionUsd := goodCoinAchievement.SelfCommissionUsd.Add(achievement.SelfCommission)
				if _, err := tx.
					GoodCoinAchievement.
					UpdateOneID(goodCoinAchievement.ID).
					SetTotalUnits(totalUnit).
					SetSelfUnits(selfUnits).
					SetTotalAmountUsd(totalAmountUsd).
					SetSelfAmountUsd(selfAmountUsd).
					SetTotalCommissionUsd(totalCommissionUsd).
					SetSelfCommissionUsd(selfCommissionUsd).
					Save(ctx); err != nil {
					return wlog.WrapError(err)
				}
				// when update exist record, we also need to migrate old one to new table, but set deleted_at = current time
				continue
			}

			if _, err := tx.
				GoodCoinAchievement.
				Create().
				SetEntID(achievement.EntID).
				SetAppID(achievement.AppID).
				SetUserID(achievement.UserID).
				SetGoodCoinTypeID(achievement.CoinTypeID).
				SetTotalUnits(achievement.TotalUnitsV1).
				SetSelfUnits(achievement.SelfUnitsV1).
				SetTotalAmountUsd(achievement.TotalAmount).
				SetSelfAmountUsd(achievement.SelfAmount).
				SetTotalCommissionUsd(achievement.TotalCommission).
				SetSelfCommissionUsd(achievement.SelfCommission).
				SetCreatedAt(achievement.CreatedAt).
				SetUpdatedAt(achievement.UpdatedAt).
				Save(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

//nolint:gocyclo
func migrateAchievementStatement(ctx context.Context, tx *ent.Tx) error {
	type OrderStatement struct {
		EntID uuid.UUID `json:"ent_id"`
	}

	rows, err := tx.QueryContext(ctx, "select ent_id from order_statements")
	if err != nil {
		return wlog.WrapError(err)
	}
	orderStatements := map[uuid.UUID]bool{}
	for rows.Next() {
		orderStatement := &OrderStatement{}
		if err := rows.Scan(&orderStatement.EntID); err != nil {
			return wlog.WrapError(err)
		}
		orderStatements[orderStatement.EntID] = true
	}

	rows, err = tx.QueryContext(ctx, "select ent_id,app_id,user_id,good_id,app_good_id,order_id,self_order,direct_contributor_id,coin_type_id,units_v1,usd_amount,app_config_id,commission_config_id,commission_config_type,payment_coin_type_id,amount,commission,payment_coin_usd_currency,created_at,updated_at from archivement_details where deleted_at = 0") //nolint
	if err != nil {
		return err
	}

	type Statement struct {
		EntID                  uuid.UUID       `json:"ent_id"`
		AppID                  uuid.UUID       `json:"app_id"`
		UserID                 uuid.UUID       `json:"user_id"`
		GoodID                 uuid.UUID       `json:"good_id"`
		AppGoodID              uuid.UUID       `json:"app_good_id"`
		OrderID                uuid.UUID       `json:"order_id"`
		SelfOrder              bool            `json:"self_order"`
		DirectContributorID    uuid.UUID       `json:"direct_contributor_id"`
		CoinTypeID             uuid.UUID       `json:"coin_type_id"`
		UnitsV1                decimal.Decimal `json:"units_v1"`
		UsdAmount              decimal.Decimal `json:"usd_amount"`
		AppConfigID            uuid.UUID       `json:"app_config_id"`
		CommissionConfigID     uuid.UUID       `json:"commission_config_id"`
		CommissionConfigType   string          `json:"commission_config_type"`
		PaymentCoinTypeID      uuid.UUID       `json:"payment_coin_type_id"`
		PaymentCoinUsdCurrency decimal.Decimal `json:"payment_coin_usd_currency"`
		Amount                 decimal.Decimal `json:"amount"`
		Commission             decimal.Decimal `json:"commission"`
		CreatedAt              uint32          `json:"created_at"`
		UpdatedAt              uint32          `json:"updated_at"`
	}

	statements := []*Statement{}
	for rows.Next() {
		statement := &Statement{}
		if err := rows.Scan(
			&statement.EntID,
			&statement.AppID,
			&statement.UserID,
			&statement.GoodID,
			&statement.AppGoodID,
			&statement.OrderID,
			&statement.SelfOrder,
			&statement.DirectContributorID,
			&statement.CoinTypeID,
			&statement.UnitsV1,
			&statement.UsdAmount,
			&statement.AppConfigID,
			&statement.CommissionConfigID,
			&statement.CommissionConfigType,
			&statement.PaymentCoinTypeID,
			&statement.Amount,
			&statement.Commission,
			&statement.PaymentCoinUsdCurrency,
			&statement.CreatedAt,
			&statement.UpdatedAt,
		); err != nil {
			return err
		}
		statements = append(statements, statement)
	}

	orderUser := map[uuid.UUID]uuid.UUID{}
	for _, statement := range statements {
		if statement.DirectContributorID != uuid.Nil {
			continue
		}
		orderUser[statement.OrderID] = statement.UserID
	}

	orderIDs := []string{}
	for _, statement := range statements {
		if statement.Amount.Cmp(decimal.NewFromInt(0)) > 0 && statement.UsdAmount.Cmp(decimal.NewFromInt(0)) == 0 {
			orderIDs = append(orderIDs, statement.OrderID.String())
		}
	}

	infos, _, err := powerrentalordermwcli.GetPowerRentalOrders(ctx, &powerrentalordermwpb.Conds{
		OrderIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: orderIDs},
	}, 0, int32(len(orderIDs)))
	if err != nil {
		return wlog.WrapError(err)
	}
	orders := map[string]*powerrentalordermwpb.PowerRentalOrder{}
	for _, info := range infos {
		orders[info.OrderID] = info
	}

	for _, statement := range statements {
		commissionAmountUSD := statement.Commission.Mul(statement.PaymentCoinUsdCurrency)
		_, ok := orderStatements[statement.EntID]
		if !ok {
			directContributorID := statement.DirectContributorID
			orderUserID, ok := orderUser[statement.OrderID]
			if !ok {
				orderUserID = uuid.Nil
			}
			if statement.UserID == orderUserID {
				directContributorID = statement.UserID
			}

			goodValueUsd := statement.UsdAmount
			paymentAmountUsd := statement.UsdAmount
			if statement.Amount.Cmp(decimal.NewFromInt(0)) > 0 && statement.UsdAmount.Cmp(decimal.NewFromInt(0)) == 0 {
				order, ok := orders[statement.OrderID.String()]
				if ok {
					switch order.OrderType {
					case ordertypes.OrderType_Offline:
						fallthrough //nolint
					case ordertypes.OrderType_Airdrop:
						goodValueUsd = statement.Amount
						paymentAmountUsd = statement.Amount
					}
				}
			}
			if _, err := tx.
				OrderStatement.
				Create().
				SetEntID(statement.EntID).
				SetAppID(statement.AppID).
				SetUserID(statement.UserID).
				SetGoodID(statement.GoodID).
				SetAppGoodID(statement.AppGoodID).
				SetOrderID(statement.OrderID).
				SetOrderUserID(orderUserID).
				SetDirectContributorID(directContributorID).
				SetGoodCoinTypeID(statement.CoinTypeID).
				SetUnits(statement.UnitsV1).
				SetGoodValueUsd(goodValueUsd).
				SetPaymentAmountUsd(paymentAmountUsd).
				SetCommissionAmountUsd(commissionAmountUSD).
				SetAppConfigID(statement.AppConfigID).
				SetCommissionConfigID(statement.CommissionConfigID).
				SetCommissionConfigType(statement.CommissionConfigType).
				SetCreatedAt(statement.CreatedAt).
				SetUpdatedAt(statement.UpdatedAt).
				Save(ctx); err != nil {
				return err
			}
		}
		exist, err := tx.
			OrderPaymentStatement.
			Query().
			Where(
				entorderpaymentstatement.StatementID(statement.EntID),
				entorderpaymentstatement.PaymentCoinTypeID(statement.PaymentCoinTypeID),
				entorderpaymentstatement.DeletedAt(0),
			).
			Exist(ctx)
		if err != nil {
			return err
		}
		if !exist {
			if _, err := tx.
				OrderPaymentStatement.
				Create().
				SetStatementID(statement.EntID).
				SetPaymentCoinTypeID(statement.PaymentCoinTypeID).
				SetAmount(statement.Amount).
				SetCommissionAmount(statement.Commission).
				SetCreatedAt(statement.CreatedAt).
				SetUpdatedAt(statement.UpdatedAt).
				Save(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func Migrate(ctx context.Context) error {
	logger.Sugar().Infow("Migrate inspire", "Start", "...")
	if err := redis2.TryLock(lockKey(), 0); err != nil {
		return wlog.WrapError(err)
	}
	defer func() {
		_ = redis2.Unlock(lockKey())
	}()

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		if err := migrateAchievementStatement(_ctx, tx); err != nil {
			logger.Sugar().Errorw("Migrate", "error", err)
			return err
		}
		if err := migrateAchievement(_ctx, tx); err != nil {
			logger.Sugar().Errorw("Migrate", "error", err)
			return err
		}
		logger.Sugar().Warnf("Migrate", "Done", "success")
		return nil
	})
}
