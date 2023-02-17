//nolint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NpoolPlatform/inspire-manager/pkg/db"
	"github.com/NpoolPlatform/inspire-manager/pkg/db/ent"
	"github.com/shopspring/decimal"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
)

const (
	keyUsername = "username"
	keyPassword = "password"
	keyDBName   = "database_name"
	maxOpen     = 10
	maxIdle     = 10
	MaxLife     = 3
)

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
		logger.Sugar().Infow("Migrate order", "Done", "...", "error", err)
	}()

	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		infos, err := cli.
			ArchivementDetail.
			Query().
			All(_ctx)
		if err != nil {
			return err
		}

		for _, info := range infos {
			if info.Units == 0 {
				continue
			}
			units := decimal.NewFromInt32(int32(info.Units))
			_, err := cli.
				ArchivementDetail.
				UpdateOneID(info.ID).
				SetUnitsV1(units).
				Save(_ctx)
			if err != nil {
				return err
			}
		}
		infos1, err := cli.
			ArchivementGeneral.
			Query().
			All(_ctx)
		if err != nil {
			return err
		}

		for _, info := range infos1 {
			u := cli.
				ArchivementGeneral.
				UpdateOneID(info.ID)

			if info.TotalUnits != 0 {
				totalUnits := decimal.NewFromInt32(int32(info.TotalUnits))
				u.SetTotalUnitsV1(totalUnits)
			}

			if info.SelfUnits != 0 {
				selfUnits := decimal.NewFromInt32(int32(info.SelfUnits))
				u.SetSelfUnitsV1(selfUnits)
			}
			_, err = u.Save(_ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
