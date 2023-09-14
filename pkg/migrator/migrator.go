//nolint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	servicename "github.com/NpoolPlatform/inspire-gateway/pkg/servicename"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db/ent"
	entachievement "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/achievement"
	entcommission "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/commission"
	entstatement "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/statement"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"

	"github.com/google/uuid"
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
	return fmt.Sprintf("%v:%v", basetypes.Prefix_PrefixMigrate, serviceID)
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

func migrateCommission(ctx context.Context, tx *ent.Tx) error {
	commissions, err := tx.
		Commission.
		Query().
		Where(
			entcommission.DeletedAt(0),
		).
		All(ctx)
	if err != nil {
		return err
	}

	goods := map[uuid.UUID]*appgoodmwpb.Good{}
	for _, commission := range commissions {
		if commission.AppGoodID != uuid.Nil {
			continue
		}
		good, ok := goods[commission.GoodID]
		if !ok {
			good, err = appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
				AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: commission.AppID.String()},
				GoodID: &basetypes.StringVal{Op: cruder.EQ, Value: commission.GoodID.String()},
			})
			if err != nil {
				continue
			}
			if good == nil {
				continue
			}
			goods[commission.GoodID] = good
		}

		id := uuid.MustParse(good.ID)
		if _, err := tx.
			Commission.
			UpdateOneID(commission.ID).
			SetAppGoodID(id).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
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

	goods := map[uuid.UUID]*appgoodmwpb.Good{}
	for _, achievement := range achievements {
		if achievement.AppGoodID != uuid.Nil {
			continue
		}
		good, ok := goods[achievement.GoodID]
		if !ok {
			good, err = appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
				AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: achievement.AppID.String()},
				GoodID: &basetypes.StringVal{Op: cruder.EQ, Value: achievement.GoodID.String()},
			})
			if err != nil {
				continue
			}
			if good == nil {
				continue
			}
			goods[achievement.GoodID] = good
		}

		id := uuid.MustParse(good.ID)
		if _, err := tx.
			Achievement.
			UpdateOneID(achievement.ID).
			SetAppGoodID(id).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func migrateStatement(ctx context.Context, tx *ent.Tx) error {
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

	goods := map[uuid.UUID]*appgoodmwpb.Good{}
	for _, statement := range statements {
		if statement.AppGoodID != uuid.Nil {
			continue
		}
		good, ok := goods[statement.GoodID]
		if !ok {
			good, err = appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
				AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: statement.AppID.String()},
				GoodID: &basetypes.StringVal{Op: cruder.EQ, Value: statement.GoodID.String()},
			})
			if err != nil {
				continue
			}
			if good == nil {
				continue
			}
			goods[statement.GoodID] = good
		}

		id := uuid.MustParse(good.ID)
		if _, err := tx.
			Statement.
			UpdateOneID(statement.ID).
			SetAppGoodID(id).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func Migrate(ctx context.Context) error {
	if err := redis2.TryLock(lockKey(), 0); err != nil {
		return err
	}
	defer func() {
		_ = redis2.Unlock(lockKey())
	}()

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		if err := migrateCommission(_ctx, tx); err != nil {
			return err
		}
		if err := migrateAchievement(_ctx, tx); err != nil {
			return err
		}
		if err := migrateStatement(_ctx, tx); err != nil {
			return err
		}
		return nil
	})
}
