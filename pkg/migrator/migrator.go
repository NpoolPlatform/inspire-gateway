//nolint
package migrator

import (
	"context"
	"fmt"
	"strings"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/mysql"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	servicename "github.com/NpoolPlatform/inspire-gateway/pkg/servicename"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db/ent"
	entstatement "github.com/NpoolPlatform/inspire-middleware/pkg/db/ent/statement"
	inspiretypes "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	ordertypes "github.com/NpoolPlatform/message/npool/basetypes/order/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/shopspring/decimal"

	"github.com/google/uuid"
)

const (
	keyServiceID = "serviceid"
)

func lockKey() string {
	serviceID := config.GetStringValueWithNameSpace(servicename.ServiceDomain, keyServiceID)
	return fmt.Sprintf("%v:%v", basetypes.Prefix_PrefixMigrate, serviceID)
}

func CreateSubordinateProcedure(ctx context.Context) error {
	conn, err := mysql.GetConn()
	if err != nil {
		return err
	}

	const procedure = `
		DROP PROCEDURE IF EXISTS get_subordinates;
		SET SESSION GROUP_CONCAT_MAX_LEN = 1024000;
		CREATE PROCEDURE get_subordinates (IN inviters TEXT)
		BEGIN
		  DECLARE subordinates TEXT;
		  DECLARE my_inviters TEXT;
		  SET subordinates = 'N/A';
		  SET my_inviters = inviters;
		  WHILE my_inviters is not null DO
		    if subordinates = 'N/A' THEN
			  SET subordinates = my_inviters;
			else
			  SET subordinates = CONCAT(subordinates, ',', my_inviters);
			END if;
		    SELECT GROUP_CONCAT(DISTINCT invitee_id) INTO my_inviters FROM registrations WHERE FIND_IN_SET(inviter_id, my_inviters) AND deleted_at=0;
		  END WHILE;
		  SELECT subordinates;
		END
	`
	_, err = conn.ExecContext(ctx, procedure)
	if err != nil {
		return err
	}

	return nil
}

func getInvites(ctx context.Context, tx *ent.Tx, inviterID string) ([]uuid.UUID, error) {
	inviterIDs := []uuid.UUID{}
	selectInviteeIDsStr := fmt.Sprintf("CALL get_subordinates(\"%v\")\n", inviterID)
	logger.Sugar().Infow("Migrate inspire", "exec selectInviteeIDsStr", selectInviteeIDsStr)
	rows, err := tx.QueryContext(
		ctx,
		selectInviteeIDsStr,
	)
	if err != nil {
		return inviterIDs, err
	}
	defer rows.Close()

	subordinates := ""
	for rows.Next() {
		if err := rows.Scan(&subordinates); err != nil {
			return []uuid.UUID{}, err
		}
	}
	_inviterIDs := strings.Split(subordinates, ",") //nolint
	for _, id := range _inviterIDs {
		if inviterID == id {
			continue
		}
		_id, err := uuid.Parse(id)
		if err != nil {
			return []uuid.UUID{}, err
		}
		inviterIDs = append(inviterIDs, _id)
	}
	return inviterIDs, nil
}

func getDirectInvites(ctx context.Context, tx *ent.Tx, userID, appID string) ([]uuid.UUID, error) {
	directInviterIDs := []uuid.UUID{}
	selectDirectInvitesStr := fmt.Sprintf("select id,app_id,inviter_id,invitee_id,deleted_at from registrations where inviter_id = '%v' and app_id='%v' and deleted_at=0", userID, appID)
	logger.Sugar().Infow("Migrate inspire", "exec selectDirectInvitesStr", selectDirectInvitesStr)
	r, err := tx.QueryContext(ctx, selectDirectInvitesStr)
	if err != nil {
		return directInviterIDs, err
	}
	type rg struct {
		ID        uint32
		AppID     uuid.UUID
		InviterID uuid.UUID
		InviteeID uuid.UUID
		DeletedAt uint32
	}
	for r.Next() {
		reg := &rg{}
		if err := r.Scan(&reg.ID, &reg.AppID, &reg.InviterID, &reg.InviteeID, &reg.DeletedAt); err != nil {
			return []uuid.UUID{}, err
		}
		directInviterIDs = append(directInviterIDs, reg.InviteeID)
	}
	r.Close()
	return directInviterIDs, nil
}

func getPaymentAmount(ctx context.Context, tx *ent.Tx, userIDs []uuid.UUID, appID uuid.UUID) (decimal.Decimal, error) {
	var err error
	paymentAmount := decimal.NewFromInt(0)
	if len(userIDs) == 0 {
		return paymentAmount, nil
	}
	var sb strings.Builder
	sb.WriteString("(")
	i := 0
	for _, id := range userIDs {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("'%v'", id))
		i++
	}
	sb.WriteString(")")
	stateStr := fmt.Sprintf("('%v', '%v', '%v')", ordertypes.OrderState_OrderStatePaid.String(), ordertypes.OrderState_OrderStateInService.String(), ordertypes.OrderState_OrderStateExpired.String())
	selectOrderStr := fmt.Sprintf("select a.id,a.app_id,a.user_id,a.payment_amount,b.order_state as state,a.deleted_at from order_manager.orders a left join order_manager.order_states b on a.ent_id=b.order_id where a.app_id='%v' and a.user_id in %s and b.order_state in %v and a.deleted_at=0", appID, sb.String(), stateStr)
	logger.Sugar().Infow("Migrate inspire", "exec selectOrderStr", selectOrderStr)
	r, err := tx.QueryContext(ctx, selectOrderStr)
	if err != nil {
		return paymentAmount, err
	}
	type od struct {
		ID            uint32
		AppID         uuid.UUID
		UserID        uuid.UUID
		PaymentAmount decimal.Decimal
		State         string
		DeletedAt     uint32
	}
	orders := []*od{}
	for r.Next() {
		order := &od{}
		if err := r.Scan(&order.ID, &order.AppID, &order.UserID, &order.PaymentAmount, &order.State, &order.DeletedAt); err != nil {
			return paymentAmount, err
		}
		orders = append(orders, order)
	}
	r.Close()

	for _, order := range orders {
		_paymentAmount := order.PaymentAmount
		if _paymentAmount.Cmp(decimal.NewFromInt(0)) > 0 {
			paymentAmount = paymentAmount.Add(_paymentAmount)
		}
	}
	return paymentAmount, nil
}

func migrateStatement(ctx context.Context, tx *ent.Tx) error {
	selectAppConfigStr := "select id, ent_id, app_id, start_at, end_at, deleted_at from app_configs where deleted_at=0 and end_at=0"
	r, err := tx.QueryContext(ctx, selectAppConfigStr)
	logger.Sugar().Warnw("Migrate inspire", "exec selectAppConfigStr", selectAppConfigStr)
	if err != nil {
		return err
	}
	type cf struct {
		ID        uint32
		EntID     uuid.UUID
		AppID     uuid.UUID
		StartAt   uint32
		EndAt     uint32
		DeletedAt uint32
	}
	configMap := map[uuid.UUID]*cf{}
	newConfigMap := map[uuid.UUID]*cf{}
	for r.Next() {
		conf := &cf{}
		if err := r.Scan(&conf.ID, &conf.EntID, &conf.AppID, &conf.StartAt, &conf.EndAt, &conf.DeletedAt); err != nil {
			return err
		}
		configMap[conf.AppID] = conf
	}

	selectAppStr := "select id,ent_id,deleted_at from appuser_manager.apps where deleted_at=0"
	r, err = tx.QueryContext(ctx, selectAppStr)
	logger.Sugar().Warnw("Migrate inspire", "exec selectAppStr", selectAppStr)
	if err != nil {
		return err
	}
	type ap struct {
		ID        uint32
		EntID     uuid.UUID
		DeletedAt uint32
	}
	for r.Next() {
		app := &ap{}
		if err := r.Scan(&app.ID, &app.EntID, &app.DeletedAt); err != nil {
			return err
		}
		_, ok := configMap[app.EntID]
		if !ok {
			id := uuid.New()
			newConf := &cf{
				EntID:   id,
				AppID:   app.EntID,
				StartAt: 0,
				EndAt:   0,
			}
			newConfigMap[app.EntID] = newConf
			configMap[app.EntID] = newConf
		}
	}

	for _, conf := range newConfigMap {
		if _, err := tx.
			AppConfig.
			Create().
			SetEntID(conf.EntID).
			SetAppID(conf.AppID).
			SetSettleMode(inspiretypes.SettleMode_SettleWithGoodValue.String()).
			SetSettleAmountType(inspiretypes.SettleAmountType_SettleByPercent.String()).
			SetSettleInterval(inspiretypes.SettleInterval_SettleEveryOrder.String()).
			SetCommissionType(inspiretypes.CommissionConfigType_LegacyCommissionConfig.String()).
			SetSettleBenefit(false).
			SetStartAt(conf.StartAt).
			SetEndAt(conf.EndAt).
			Save(ctx); err != nil {
			return err
		}
	}

	offset := 0
	limit := 1000
	for {
		selectStatementStr := fmt.Sprintf("select id,app_id,ent_id,deleted_at from archivement_details where deleted_at=0 and commission_config_type='%v' limit %v, %v", inspiretypes.CommissionConfigType_DefaultCommissionConfigType.String(), offset, limit)
		logger.Sugar().Warnw("Migrate inspire", "exec selectStatementStr", selectStatementStr)
		r, err := tx.QueryContext(ctx, selectStatementStr)
		if err != nil {
			return err
		}
		type ad struct {
			ID        uint32
			AppID     uuid.UUID
			EntID     uuid.UUID
			DeletedAt uint32
		}
		statements := []*ad{}
		for r.Next() {
			statement := &ad{}
			if err := r.Scan(&statement.ID, &statement.AppID, &statement.EntID, &statement.DeletedAt); err != nil {
				return err
			}
			statements = append(statements, statement)
		}

		logger.Sugar().Warnw("Migrate inspire", "exec len(statements)", len(statements))
		if len(statements) == 0 {
			break
		}

		for _, statement := range statements {
			appConfigID := uuid.Nil
			conf, ok := configMap[statement.AppID]
			if ok {
				appConfigID = conf.EntID
			}
			if _, err := tx.
				Statement.
				Update().
				Where(
					entstatement.ID(statement.ID),
				).
				SetAppConfigID(appConfigID).
				SetCommissionConfigID(uuid.Nil).
				SetCommissionConfigType(inspiretypes.CommissionConfigType_LegacyCommissionConfig.String()).
				Save(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func migrateAchievementUser(ctx context.Context, tx *ent.Tx) error {
	var err error
	r, err := tx.QueryContext(ctx, "select id,app_id,ent_id,deleted_at from appuser_manager.app_users where deleted_at=0")
	if err != nil {
		return err
	}
	type us struct {
		ID        uint32
		AppID     uuid.UUID
		EntID     uuid.UUID
		DeletedAt uint32
	}
	users := []*us{}
	for r.Next() {
		user := &us{}
		if err := r.Scan(&user.ID, &user.AppID, &user.EntID, &user.DeletedAt); err != nil {
			return err
		}
		users = append(users, user)
	}
	r.Close()
	if len(users) == 0 {
		return nil
	}

	logger.Sugar().Warnw("Migrate inspire", "exec len(users)", len(users))
	count := 0
	for _, user := range users {
		count++
		selectAchievementUserStr := fmt.Sprintf("select id,app_id,user_id,deleted_at from archivement_users where user_id='%v' and deleted_at=0", user.EntID)
		logger.Sugar().Warnw("Migrate inspire", "exec selectAchievementUserStr", selectAchievementUserStr)
		r, err = tx.QueryContext(ctx, selectAchievementUserStr)
		if err != nil {
			return err
		}
		type achievedUser struct {
			ID                   uint32
			EntID                uuid.UUID
			AppID                uuid.UUID
			UserID               uuid.UUID
			TotalCommission      decimal.Decimal
			SelfCommission       decimal.Decimal
			DirectInvites        uint32
			IndirectInvites      uint32
			DirectConsumeAmount  decimal.Decimal
			InviteeConsumeAmount decimal.Decimal
			CreatedAt            uint32
			UpdatedAt            uint32
			DeletedAt            uint32
		}
		achievedUsers := []*achievedUser{}
		for r.Next() {
			achievementUser := &achievedUser{}
			if err := r.Scan(&achievementUser.ID, &achievementUser.AppID, &achievementUser.UserID, &achievementUser.DeletedAt); err != nil {
				return err
			}
			achievedUsers = append(achievedUsers, achievementUser)
		}
		r.Close()
		if len(achievedUsers) != 0 {
			continue
		}

		id := uuid.New()
		_achievedUser := &achievedUser{
			EntID:  id,
			AppID:  user.AppID,
			UserID: user.EntID,
		}

		// reigster
		inviteeIDs, err := getInvites(ctx, tx, user.EntID.String())
		if err != nil {
			return err
		}
		directInviteeIDs, err := getDirectInvites(ctx, tx, user.EntID.String(), user.AppID.String())
		if err != nil {
			return err
		}

		_achievedUser.DirectInvites = uint32(len(directInviteeIDs))
		_achievedUser.IndirectInvites = uint32(len(inviteeIDs) - len(directInviteeIDs))

		// order
		userIDs := []uuid.UUID{}
		userIDs = append(userIDs, user.EntID)
		directConsumeAmount, err := getPaymentAmount(ctx, tx, userIDs, user.AppID)
		if err != nil {
			return err
		}
		indirectConsumeAmount, err := getPaymentAmount(ctx, tx, inviteeIDs, user.AppID)
		if err != nil {
			return err
		}

		_achievedUser.DirectConsumeAmount = directConsumeAmount
		_achievedUser.InviteeConsumeAmount = indirectConsumeAmount

		// commission
		selectAchievementStr := fmt.Sprintf("select id,app_id,user_id,total_commission,self_commission,deleted_at from archivement_generals where app_id='%v' and user_id='%v' and deleted_at=0", user.AppID, user.EntID)
		logger.Sugar().Warnw("Migrate inspire", "exec selectAchievementStr", selectAchievementStr)
		r, err = tx.QueryContext(ctx, selectAchievementStr)
		if err != nil {
			return err
		}
		type achv struct {
			ID              uint32
			AppID           uuid.UUID
			UserID          uuid.UUID
			TotalCommission decimal.Decimal
			SelfCommission  decimal.Decimal
			DeletedAt       uint32
		}
		achievements := []*achv{}
		for r.Next() {
			achievement := &achv{}
			if err := r.Scan(&achievement.ID, &achievement.AppID, &achievement.UserID, &achievement.TotalCommission, &achievement.SelfCommission, &achievement.DeletedAt); err != nil {
				return err
			}
			achievements = append(achievements, achievement)
		}
		r.Close()

		totalCommission := decimal.NewFromInt(0)
		selfCommission := decimal.NewFromInt(0)
		for _, achievement := range achievements {
			totalCommission = totalCommission.Add(achievement.TotalCommission)
			selfCommission = selfCommission.Add(achievement.SelfCommission)
		}

		_achievedUser.TotalCommission = totalCommission
		_achievedUser.SelfCommission = selfCommission

		logger.Sugar().Warnw("Migrate inspire", "create _achievedUser", _achievedUser)
		if _, err := tx.
			AchievementUser.
			Create().
			SetEntID(_achievedUser.EntID).
			SetAppID(_achievedUser.AppID).
			SetUserID(_achievedUser.UserID).
			SetTotalCommission(_achievedUser.TotalCommission).
			SetSelfCommission(_achievedUser.SelfCommission).
			SetDirectInvites(_achievedUser.DirectInvites).
			SetIndirectInvites(_achievedUser.IndirectInvites).
			SetDirectConsumeAmount(_achievedUser.DirectConsumeAmount).
			SetInviteeConsumeAmount(_achievedUser.InviteeConsumeAmount).
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

	if err := CreateSubordinateProcedure(ctx); err != nil {
		return err
	}

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		if err := migrateAchievementUser(_ctx, tx); err != nil {
			logger.Sugar().Errorw("Migrate", "error", err)
			return err
		}
		if err := migrateStatement(_ctx, tx); err != nil {
			logger.Sugar().Errorw("Migrate", "error", err)
			return err
		}
		logger.Sugar().Infow("Migrate", "Done", "success")
		return nil
	})
}
