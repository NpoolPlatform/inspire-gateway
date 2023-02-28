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
		details, err := tx.
			ArchivementDetail.
			Query().
			Where(
				archivementdetailent.DeletedAt(0),
			).
			All(_ctx)
		if err != nil {
			return err
		}

		type order struct {
			AppID   uuid.UUID
			UserID  uuid.UUID
			OrderID uuid.UUID
		}

		rows, err := tx.
			QueryContext(
				_ctx,
				"select app_id,user_id,id from order_manager.orders where type='Normal' and deleted_at=0")
		if err != nil {
			return err
		}

		orders := []*order{}
		for rows.Next() {
			order := &order{}
			err := rows.Scan(
				&order.AppID,
				&order.UserID,
				&order.OrderID,
			)
			if err != nil {
				return err
			}
			orders = append(orders, order)
		}

		orderMap := map[string]*order{}
		for i, ord := range orders {
			orderMap[ord.OrderID.String()] = orders[i]
		}

		comms := map[string]map[string]map[string]decimal.Decimal{}
		selfComms := map[string]map[string]map[string]decimal.Decimal{}

		for _, info := range details {
			ord, ok := orderMap[info.OrderID.String()]
			if !ok {
				if info.Commission.Cmp(decimal.NewFromInt(0)) > 0 {
					logger.Sugar().Infow("Migrate", "OrderID", info.OrderID, "State", "Offline | Airdrop", "Commission", info.Commission)
				}
			}

			acomm, ok := comms[info.AppID.String()]
			if !ok {
				acomm = map[string]map[string]decimal.Decimal{}
			}

			comm, ok := acomm[info.UserID.String()]
			if !ok {
				comm = map[string]decimal.Decimal{}
			}

			_comm, ok := comm[info.GoodID.String()]
			if !ok {
				_comm = decimal.Decimal{}
			}

			_comm = _comm.Add(info.Commission)
			comm[info.GoodID.String()] = _comm
			acomm[info.UserID.String()] = comm

			comms[info.AppID.String()] = acomm

			if ord == nil {
				continue
			}

			logger.Sugar().Infow("Migrate", "OrderID", info.OrderID, "Comm", _comm)
			if info.UserID != ord.UserID {
				continue
			}

			acomm, ok = selfComms[info.AppID.String()]
			if !ok {
				acomm = map[string]map[string]decimal.Decimal{}
			}

			comm, ok = acomm[info.UserID.String()]
			if !ok {
				comm = map[string]decimal.Decimal{}
			}

			_comm, ok = comm[info.GoodID.String()]
			if !ok {
				_comm = decimal.Decimal{}
			}

			_comm = _comm.Add(info.Commission)
			logger.Sugar().Infow("Migrate", "OrderID", info.OrderID, "SelfComm", _comm)
			comm[info.GoodID.String()] = _comm
			acomm[info.UserID.String()] = comm

			selfComms[info.AppID.String()] = acomm
		}

		for appID, _comms := range comms {
			for userID, __comms := range _comms {
				for goodID, comm := range __comms {
					logger.Sugar().Infow("Migrate", "AppID", appID, "UserID", userID, "GoodID", goodID)
					general, err := tx.
						ArchivementGeneral.
						Query().
						Where(
							archivementgeneralent.AppID(uuid.MustParse(appID)),
							archivementgeneralent.UserID(uuid.MustParse(userID)),
							archivementgeneralent.GoodID(uuid.MustParse(goodID)),
						).
						Only(_ctx)
					if err != nil {
						if ent.IsNotFound(err) {
							logger.Sugar().Errorw(
								"Migrate",
								"AppID", appID,
								"UserID", userID,
								"GoodID", goodID,
								"Error", "Invalid General")
							continue
						}
						return err
					}

					selfComm := decimal.Decimal{}

					_selfComms, ok := selfComms[appID]
					if ok {
						__selfComms, ok := _selfComms[userID]
						if ok {
							_selfComm, ok := __selfComms[goodID]
							if ok {
								selfComm = _selfComm
							}
						}
					}

					_, err = tx.
						ArchivementGeneral.
						UpdateOneID(general.ID).
						SetTotalCommission(comm).
						SetSelfCommission(selfComm).
						Save(_ctx)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	return err
}
