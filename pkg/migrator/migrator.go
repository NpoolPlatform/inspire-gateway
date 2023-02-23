//nolint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NpoolPlatform/go-service-framework/pkg/redis"
	"github.com/NpoolPlatform/inspire-manager/pkg/db"
	"github.com/NpoolPlatform/inspire-manager/pkg/db/ent"
	archivementdetailent "github.com/NpoolPlatform/inspire-manager/pkg/db/ent/archivementdetail"
	archivementgeneralent "github.com/NpoolPlatform/inspire-manager/pkg/db/ent/archivementgeneral"
	"github.com/shopspring/decimal"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	constant1 "github.com/NpoolPlatform/inspire-gateway/pkg/message/const"
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
	serviceID := config.GetStringValueWithNameSpace(constant1.ServiceName, keyServiceID)
	return fmt.Sprintf("migrator:%v", serviceID)
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
	var err error

	if err := db.Init(); err != nil {
		return err
	}
	logger.Sugar().Infow("Migrate order", "Start", "...")
	defer func() {
		_ = redis.Unlock(lockKey())
		logger.Sugar().Infow("Migrate order", "Done", "...", "error", err)
	}()

	err = redis.TryLock(lockKey(), 0)
	if err != nil {
		return err
	}

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		infos, err := tx.
			ArchivementDetail.
			Query().
			Select(
				archivementdetailent.FieldID,
				archivementdetailent.FieldUnits,
			).
			All(_ctx)
		if err != nil {
			return err
		}

		for _, info := range infos {
			units := decimal.NewFromInt(0)
			if info.Units != 0 {
				units = decimal.NewFromInt32(int32(info.Units))
			}

			_, err = tx.
				ArchivementDetail.
				UpdateOneID(info.ID).
				SetUnitsV1(units).
				Save(_ctx)
			if err != nil {
				return err
			}
		}
		infos1, err := tx.
			ArchivementGeneral.
			Query().
			Select(
				archivementgeneralent.FieldID,
				archivementgeneralent.FieldTotalUnits,
				archivementgeneralent.FieldSelfUnits,
			).
			All(_ctx)
		if err != nil {
			return err
		}

		for _, info := range infos1 {
			u := tx.
				ArchivementGeneral.
				UpdateOneID(info.ID)

			totalUnits := decimal.NewFromInt(0)
			if info.TotalUnits != 0 {
				totalUnits = decimal.NewFromInt32(int32(info.TotalUnits))
			}
			u.SetTotalUnitsV1(totalUnits)

			selfUnits := decimal.NewFromInt(0)
			if info.SelfUnits != 0 {
				selfUnits = decimal.NewFromInt32(int32(info.SelfUnits))
			}
			u.SetSelfUnitsV1(selfUnits)
			_, err = u.Save(_ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
