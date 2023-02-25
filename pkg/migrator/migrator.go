//nolint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/redis"
	"github.com/NpoolPlatform/inspire-manager/pkg/db"
	"github.com/NpoolPlatform/inspire-manager/pkg/db/ent"

	archivementdetailent "github.com/NpoolPlatform/inspire-manager/pkg/db/ent/archivementdetail"
	archivementgeneralent "github.com/NpoolPlatform/inspire-manager/pkg/db/ent/archivementgeneral"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	ordermgrpb "github.com/NpoolPlatform/message/npool/order/mgr/v1/order"

	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	constant1 "github.com/NpoolPlatform/inspire-gateway/pkg/message/const"

	"github.com/google/uuid"
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

	err = db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
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

		type order struct {
			ID     uuid.UUID
			AppID  uuid.UUID
			UserID uuid.UUID
			State  string
			Type   string
		}

		rows, err := tx.
			QueryContext(
				ctx,
				"select "+
					"id,"+
					"app_id,"+
					"user_id,"+
					"state,"+
					"type "+
					"from order_manager.orders "+
					"where deleted_at=0",
			)
		if err != nil {
			return err
		}

		ords := []*order{}

		for rows.Next() {
			order := order{}
			err := rows.Scan(
				&order.ID,
				&order.AppID,
				&order.UserID,
				&order.State,
				&order.Type,
			)
			if err != nil {
				return err
			}

			if order.State == ordermgrpb.OrderType_Normal.String() {
				continue
			}

			ords = append(ords, &order)
		}

		for _, order := range ords {
			infos, err := tx.
				ArchivementDetail.
				Query().
				Where(
					archivementdetailent.OrderID(order.ID),
				).
				All(_ctx)
			if err != nil {
				return err
			}

			for _, info := range infos {
				if info.Commission.Cmp(decimal.NewFromInt(0)) <= 0 {
					continue
				}

				_, err := tx.
					ArchivementDetail.
					UpdateOneID(info.ID).
					SetCommission(decimal.NewFromInt(0)).
					Save(_ctx)
				if err != nil {
					return err
				}

				general, err := tx.
					ArchivementGeneral.
					Query().
					Where(
						archivementgeneralent.AppID(order.AppID),
						archivementgeneralent.AppID(order.UserID),
					).
					Only(_ctx)
				if err != nil {
					if ent.IsNotFound(err) {
						logger.Sugar().Errorw("Migrate", "AppID", order.AppID, "UserID", order.UserID, "Error", err)
						continue
					}
					return err
				}

				totalCommission := general.TotalCommission.Sub(info.Commission)
				selfCommission := general.SelfCommission

				if info.UserID == order.UserID {
					selfCommission = selfCommission.Sub(info.Commission)
				}

				_, err = tx.
					ArchivementGeneral.
					UpdateOneID(general.ID).
					SetTotalCommission(totalCommission).
					SetSelfCommission(selfCommission).
					Save(_ctx)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}
