//nolint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db/ent"
	entcommission "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/commission"
	entgoodorderpercent "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/goodorderpercent"

	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	servicename "github.com/NpoolPlatform/inspire-gateway/pkg/servicename"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/shopspring/decimal"
)

const (
	keyUsername  = "username"
	keyPassword  = "password"
	keyDBName    = "database_name"
	maxOpen      = 10
	maxIdle      = 10
	MaxLife      = 3
	keyServiceID = "serviceid"
)

func lockKey() string {
	serviceID := config.GetStringValueWithNameSpace(servicename.ServiceDomain, keyServiceID)
	return fmt.Sprintf("%v:%v", basetypes.Prefix_PrefixMigrateInspire, serviceID)
}

func dsn(hostname string) (string, error) {
	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	dbname := config.GetStringValueWithNameSpace(hostname, keyDBName)

	svc, err := config.PeekService(constant.MysqlServiceName)
	if err != nil {
		logger.Sugar().Warnw("dsn", "error", err)
		return "", err
	}

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true",
		username, password,
		svc.Address,
		svc.Port,
		dbname,
	), nil
}

func open(hostname string) (conn *sql.DB, err error) {
	hdsn, err := dsn(hostname)
	if err != nil {
		return nil, err
	}

	logger.Sugar().Infow("open", "hdsn", hdsn)

	conn, err = sql.Open("mysql", hdsn)
	if err != nil {
		return nil, err
	}

	// https://github.com/go-sql-driver/mysql
	// See "Important settings" section.

	conn.SetConnMaxLifetime(time.Minute * MaxLife)
	conn.SetMaxOpenConns(maxOpen)
	conn.SetMaxIdleConns(maxIdle)

	return conn, nil
}

func Migrate(ctx context.Context) error {
	if err := redis2.TryLock(lockKey(), 0); err != nil {
		return err
	}
	defer func() {
		_ = redis2.Unlock(lockKey())
	}()

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		gops, err := tx.
			GoodOrderPercent.
			Query().
			Where(
				entgoodorderpercent.DeletedAt(0),
			).
			All(_ctx)
		if err != nil {
			return err
		}

		for _, gop := range gops {
			exist, err := tx.
				Commission.
				Query().
				Where(
					entcommission.ID(gop.ID),
				).
				Exist(_ctx)
			if err != nil {
				return err
			}
			if exist {
				continue
			}

			if _, err := tx.
				Commission.
				Create().
				SetID(gop.ID).
				SetAppID(gop.AppID).
				SetUserID(gop.UserID).
				SetGoodID(gop.GoodID).
				SetAmountOrPercent(gop.Percent).
				SetStartAt(gop.StartAt).
				SetEndAt(gop.EndAt).
				SetSettleType(types.SettleType_GoodOrderPayment.String()).
				SetSettleMode(types.SettleMode_SettleWithGoodValue.String()).
				SetSettleAmountType(types.SettleAmountType_SettleByPercent.String()).
				SetSettleInterval(types.SettleInterval_SettleEveryOrder.String()).
				SetThreshold(decimal.NewFromInt(0)).
				SetOrderLimit(0).
				SetCreatedAt(gop.CreatedAt).
				SetUpdatedAt(gop.UpdatedAt).
				SetDeletedAt(gop.DeletedAt).
				Save(_ctx); err != nil {
				return err
			}
		}
		return nil

		return nil
	})
}
