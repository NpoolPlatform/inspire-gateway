//nolint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NpoolPlatform/inspire-manager/pkg/db"
	"github.com/NpoolPlatform/inspire-manager/pkg/db/ent"

	"github.com/shopspring/decimal"

	entarchivementdetail "github.com/NpoolPlatform/inspire-manager/pkg/db/ent/archivementdetail"
	entarchivementgeneral "github.com/NpoolPlatform/inspire-manager/pkg/db/ent/archivementgeneral"
	entgoodorderpercent "github.com/NpoolPlatform/inspire-manager/pkg/db/ent/goodorderpercent"
	entivcode "github.com/NpoolPlatform/inspire-manager/pkg/db/ent/invitationcode"
	entreg "github.com/NpoolPlatform/inspire-manager/pkg/db/ent/registration"

	entinspiresetting "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/db/ent/apppurchaseamountsetting"
	entregiv "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/db/ent/registrationinvitation"
	entoivcode "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/db/ent/userinvitationcode"

	archivementent "github.com/NpoolPlatform/archivement-manager/pkg/db/ent"
	entdetail "github.com/NpoolPlatform/archivement-manager/pkg/db/ent/detail"
	entgeneral "github.com/NpoolPlatform/archivement-manager/pkg/db/ent/general"

	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	constant1 "github.com/NpoolPlatform/inspire-gateway/pkg/message/const"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/redis"

	inspireent "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/db/ent"
	inspireconst "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/message/const"

	archivementmgrconst "github.com/NpoolPlatform/archivement-manager/pkg/message/const"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/NpoolPlatform/archivement-manager/pkg/db/ent/runtime"
	_ "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/db/ent/runtime"
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
		logger.Sugar().Warnw("dsb", "error", err)
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

func migrateInvitationCode(ctx context.Context, conn *sql.DB) error {
	cli := inspireent.NewClient(inspireent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	ivcodes, err := cli.
		UserInvitationCode.
		Query().
		Where(
			entoivcode.DeleteAt(0),
		).
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateInvitationCode", "Error", err)
		return err
	}

	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		infos, err := cli.
			InvitationCode.
			Query().
			Limit(1).
			All(_ctx)
		if err != nil {
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, code := range ivcodes {
			info, err := cli.
				InvitationCode.
				Query().
				Where(
					entivcode.ID(code.ID),
				).
				Only(_ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return err
				}
			}
			if info != nil {
				continue
			}

			_, err = cli.
				InvitationCode.
				Create().
				SetID(code.ID).
				SetAppID(code.AppID).
				SetUserID(code.UserID).
				SetInvitationCode(code.InvitationCode).
				SetConfirmed(code.Confirmed).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func migrateRegistration(ctx context.Context, conn *sql.DB) error {
	cli := inspireent.NewClient(inspireent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	regs, err := cli.
		RegistrationInvitation.
		Query().
		Where(
			entregiv.DeleteAt(0),
		).
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateRegistration", "Error", err)
		return err
	}

	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		infos, err := cli.
			Registration.
			Query().
			Limit(1).
			All(_ctx)
		if err != nil {
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, reg := range regs {
			info, err := cli.
				Registration.
				Query().
				Where(
					entreg.ID(reg.ID),
				).
				Only(_ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return err
				}
			}
			if info != nil {
				continue
			}

			_, err = cli.
				Registration.
				Create().
				SetID(reg.ID).
				SetAppID(reg.AppID).
				SetInviterID(reg.InviterID).
				SetInviteeID(reg.InviteeID).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func migrateAmountSetting(ctx context.Context, conn *sql.DB) error {
	cli := inspireent.NewClient(inspireent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	settings, err := cli.
		AppPurchaseAmountSetting.
		Query().
		Where(
			entinspiresetting.DeleteAt(0),
		).
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateAmountSetting", "Error", err)
		return err
	}

	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		infos, err := cli.
			GoodOrderPercent.
			Query().
			Limit(1).
			All(_ctx)
		if err != nil {
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, setting := range settings {
			info, err := cli.
				GoodOrderPercent.
				Query().
				Where(
					entgoodorderpercent.ID(setting.ID),
				).
				Only(_ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return err
				}
			}
			if info != nil {
				continue
			}

			percent := decimal.NewFromInt(int64(setting.Percent))

			_, err = cli.
				GoodOrderPercent.
				Create().
				SetID(setting.ID).
				SetAppID(setting.AppID).
				SetUserID(setting.UserID).
				SetUserID(setting.GoodID).
				SetPercent(percent).
				SetStartAt(setting.Start).
				SetEndAt(setting.End).
				SetCreatedAt(setting.CreateAt).
				SetUpdatedAt(setting.UpdateAt).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func migrateArchivementGeneral(ctx context.Context, conn *sql.DB) error {
	cli := archivementent.NewClient(archivementent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	generals, err := cli.
		General.
		Query().
		Where(
			entgeneral.DeletedAt(0),
		).
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateArchivementGeneral", "Error", err)
		return err
	}

	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		infos, err := cli.
			ArchivementGeneral.
			Query().
			Limit(1).
			All(_ctx)
		if err != nil {
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, general := range generals {
			info, err := cli.
				ArchivementGeneral.
				Query().
				Where(
					entarchivementgeneral.ID(general.ID),
				).
				Only(_ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return err
				}
			}
			if info != nil {
				continue
			}

			_, err = cli.
				ArchivementGeneral.
				Create().
				SetID(general.ID).
				SetAppID(general.AppID).
				SetUserID(general.UserID).
				SetUserID(general.GoodID).
				SetCoinTypeID(general.CoinTypeID).
				SetTotalUnits(general.TotalUnits).
				SetSelfUnits(general.SelfUnits).
				SetTotalAmount(general.TotalAmount).
				SetSelfAmount(general.SelfAmount).
				SetTotalCommission(general.TotalCommission).
				SetSelfCommission(general.SelfCommission).
				SetCreatedAt(general.CreatedAt).
				SetUpdatedAt(general.UpdatedAt).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func migrateArchivementDetail(ctx context.Context, conn *sql.DB) error {
	cli := archivementent.NewClient(archivementent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	details, err := cli.
		Detail.
		Query().
		Where(
			entdetail.DeletedAt(0),
		).
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateAmountSetting", "Error", err)
		return err
	}

	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		infos, err := cli.
			ArchivementDetail.
			Query().
			Limit(1).
			All(_ctx)
		if err != nil {
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, detail := range details {
			info, err := cli.
				ArchivementDetail.
				Query().
				Where(
					entarchivementdetail.ID(detail.ID),
				).
				Only(_ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return err
				}
			}
			if info != nil {
				continue
			}

			_, err = cli.
				ArchivementDetail.
				Create().
				SetID(detail.ID).
				SetAppID(detail.AppID).
				SetUserID(detail.UserID).
				SetDirectContributorID(detail.DirectContributorID).
				SetUserID(detail.GoodID).
				SetOrderID(detail.OrderID).
				SetSelfOrder(detail.SelfOrder).
				SetPaymentID(detail.PaymentID).
				SetCoinTypeID(detail.CoinTypeID).
				SetPaymentCoinTypeID(detail.PaymentCoinTypeID).
				SetPaymentCoinUsdCurrency(detail.PaymentCoinUsdCurrency).
				SetUnits(detail.Units).
				SetAmount(detail.Amount).
				SetUsdAmount(detail.UsdAmount).
				SetCommission(detail.Commission).
				SetCreatedAt(detail.CreatedAt).
				SetUpdatedAt(detail.UpdatedAt).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func Migrate(ctx context.Context) error {
	if err := db.Init(); err != nil {
		return err
	}

	serviceID := config.GetStringValueWithNameSpace(constant1.ServiceName, config.KeyServiceID)
	if err := redis.TryLock(serviceID, 0); err != nil {
		return nil
	}
	defer func() {
		_ = redis.Unlock(serviceID)
	}()

	conn, err := open(inspireconst.ServiceName)
	if err != nil {
		return err
	}

	if err := migrateInvitationCode(ctx, conn); err != nil {
		return err
	}

	if err := migrateRegistration(ctx, conn); err != nil {
		return err
	}

	if err := migrateAmountSetting(ctx, conn); err != nil {
		return err
	}

	conn, err = open(archivementmgrconst.ServiceName)
	if err != nil {
		return err
	}

	if err := migrateArchivementGeneral(ctx, conn); err != nil {
		return err
	}

	if err := migrateArchivementDetail(ctx, conn); err != nil {
		return err
	}

	return nil
}