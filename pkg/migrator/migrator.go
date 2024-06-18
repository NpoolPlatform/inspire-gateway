//nolint
package migrator

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	servicename "github.com/NpoolPlatform/inspire-gateway/pkg/servicename"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db/ent"
	entgoodachievement "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/goodachievement"
	entgoodcoinachievement "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/goodcoinachievement"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
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

func migrateAchievement(ctx context.Context, tx *ent.Tx) error {
	rows, err := tx.QueryContext(ctx, "select ent_id,app_id,user_id,good_id,app_good_id,coin_type_id,total_units_v1,self_units_v1,total_amount,self_amount,total_commission,self_commission,created_at,updated_at from archivement_generals where deleted_at = 0")
	if err != nil {
		return err
	}

	type Achievement struct {
		EntID           uuid.UUID
		AppID           uuid.UUID
		UserID          uuid.UUID
		GoodID          uuid.UUID
		AppGoodID       uuid.UUID
		CoinTypeID      uuid.UUID
		TotalUnitsV1    decimal.Decimal
		SelfUnitsV1     decimal.Decimal
		TotalAmount     decimal.Decimal
		SelfAmount      decimal.Decimal
		TotalCommission decimal.Decimal
		SelfCommission  decimal.Decimal
		CreatedAt       uint32
		UpdatedAt       uint32
	}
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
		exist, err := tx.
			GoodAchievement.
			Query().
			Where(
				entgoodachievement.AppID(achievement.AppID),
				entgoodachievement.UserID(achievement.UserID),
				entgoodachievement.GoodID(achievement.GoodID),
				entgoodachievement.AppGoodID(achievement.AppGoodID),
				entgoodachievement.DeletedAt(0),
			).Exist(ctx)
		if err != nil {
			return err
		}
		if !exist {
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
		exist, err = tx.
			GoodCoinAchievement.
			Query().
			Where(
				entgoodcoinachievement.AppID(achievement.AppID),
				entgoodcoinachievement.UserID(achievement.UserID),
				entgoodcoinachievement.GoodCoinTypeID(achievement.CoinTypeID),
				entgoodcoinachievement.DeletedAt(0),
			).Exist(ctx)
		if err != nil {
			return err
		}
		if !exist {
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

func migrateAchievementStatement(ctx context.Context, tx *ent.Tx) error {
	rows, err := tx.QueryContext(ctx, "select ent_id,app_id,user_id,good_id,app_good_id,order_id,direct_contributor_id,coin_type_id,units_v1,usd_amount,app_config_id,commission_config_id,commission_config_type,payment_coin_type_id,amount,commission,payment_coin_usd_currency,created_at,updated_at from archivement_details where deleted_at = 0")
	if err != nil {
		return err
	}

	type Statement struct {
		EntID                  uuid.UUID
		AppID                  uuid.UUID
		UserID                 uuid.UUID
		GoodID                 uuid.UUID
		AppGoodID              uuid.UUID
		OrderID                uuid.UUID
		DirectContributorID    uuid.UUID
		CoinTypeID             uuid.UUID
		UnitsV1                decimal.Decimal
		UsdAmount              decimal.Decimal
		AppConfigID            uuid.UUID
		CommissionConfigID     uuid.UUID
		CommissionConfigType   string
		PaymentCoinTypeID      uuid.UUID
		PaymentCoinUsdCurrency decimal.Decimal
		Amount                 decimal.Decimal
		Commission             decimal.Decimal
		CreatedAt              uint32
		UpdatedAt              uint32
	}

	for rows.Next() {
		statement := &Statement{}
		if err := rows.Scan(
			&statement.EntID,
			&statement.AppID,
			&statement.UserID,
			&statement.DirectContributorID,
			&statement.GoodID,
			&statement.OrderID,
			&statement.AppGoodID,
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
		if statement.PaymentCoinUsdCurrency.Cmp(decimal.NewFromInt(0)) <= 0 {
			logger.Sugar().Warn("invalid payment coin usd currency: %v", statement.PaymentCoinUsdCurrency.String())
			continue
		}
		commissionAmountUSD := statement.Commission.Mul(statement.PaymentCoinUsdCurrency)
		if _, err := tx.
			OrderStatement.
			Create().
			SetEntID(statement.EntID).
			SetAppID(statement.AppID).
			SetUserID(statement.UserID).
			SetAppGoodID(statement.AppGoodID).
			SetOrderID(statement.OrderID).
			SetOrderUserID(statement.DirectContributorID).
			SetGoodCoinTypeID(statement.CoinTypeID).
			SetUnits(statement.UnitsV1).
			SetGoodValueUsd(statement.UsdAmount).
			SetPaymentAmountUsd(statement.UsdAmount).
			SetCommissionAmountUsd(commissionAmountUSD).
			SetAppConfigID(statement.AppConfigID).
			SetCommissionConfigID(statement.CommissionConfigID).
			SetCommissionConfigType(statement.CommissionConfigType).
			Save(ctx); err != nil {
			return err
		}
		if _, err := tx.
			OrderPaymentStatement.
			Create().
			SetStatementID(statement.EntID).
			SetPaymentCoinTypeID(statement.PaymentCoinTypeID).
			SetAmount(statement.Amount).
			SetCommissionAmount(statement.Commission).
			Save(ctx); err != nil {
			return err
		}
	}

	return nil
}

func Migrate(ctx context.Context) error {
	logger.Sugar().Infow("Migrate inspire", "Start", "...")
	if err := redis2.TryLock(lockKey(), 0); err != nil {
		return err
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
		logger.Sugar().Infow("Migrate", "Done", "success")
		return nil
	})
}
