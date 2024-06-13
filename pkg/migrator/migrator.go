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
	entachievement "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/achievement"
	entstatement "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/statement"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
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
	achievements, err := tx.
		Achievement.
		Query().
		Where(
			entachievement.DeletedAt(0),
		).
		All(ctx)
	if err != nil {
		return err
	}

	for _, achievement := range achievements {
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
			Save(ctx); err != nil {
			return err
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
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func migrateAchievementStatement(ctx context.Context, tx *ent.Tx) error {
	statements, err := tx.
		Statement.
		Query().
		Where(
			entstatement.DeletedAt(0),
		).
		All(ctx)
	if err != nil {
		return err
	}

	for _, statement := range statements {
		if statement.PaymentCoinUsdCurrency.Cmp(decimal.NewFromInt(0)) <= 0 {
			logger.Sugar().Warn("invalid payment coin usd currency: %v", statement.PaymentCoinUsdCurrency.String())
			continue
		}
		paymentAmountUSD := statement.PaymentCoinUsdCurrency.Mul(statement.Amount)

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
			SetPaymentAmountUsd(paymentAmountUSD).
			SetCommissionAmountUsd(statement.Commission).
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
