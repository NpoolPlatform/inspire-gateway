package achievement

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	achievementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement"
	statementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement/statement"
	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	registrationmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/achievement"
	achievementmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement"
	statementmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement/statement"
	commissionmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
	registrationmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"

	"github.com/shopspring/decimal"
)

type queryHandler struct {
	*Handler
	registrations map[string]*registrationmwpb.Registration
	inviteIDs     []string
	achievements  map[string]*achievementmwpb.Achievement
	inviteesCount map[string]uint32
	coins         map[string]*appcoinmwpb.Coin
	users         map[string]*usermwpb.User
	goods         map[string]*appgoodmwpb.Good
	commissions   map[string]map[string]*commissionmwpb.Commission
	total         uint32
	achievedGoods map[string]map[string]struct{}
	statements    []*statementmwpb.Statement
	infoMap       map[string]*npool.Achievement
	infos         []*npool.Achievement
}

func (h *queryHandler) getInvitees(ctx context.Context) error {
	if h.Limit < int32(len(h.inviteIDs)) {
		return fmt.Errorf("limit should be greater than userids")
	}

	resetInviteIDs := false
	if h.Offset == 0 {
		h.Limit -= int32(len(h.inviteIDs))
	} else if h.Offset < int32(len(h.inviteIDs)) {
		return fmt.Errorf("offset should be greater than userids")
	} else {
		h.Offset -= int32(len(h.inviteIDs))
		resetInviteIDs = true
	}

	registrations, total, err := registrationmwcli.GetRegistrations(ctx, &registrationmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		InviterIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
	}, h.Offset, h.Limit)
	if err != nil {
		return err
	}
	h.total = total + uint32(len(h.inviteIDs))
	if resetInviteIDs {
		h.inviteIDs = []string{}
	}
	for i, registration := range registrations {
		if int32(i) >= h.Limit {
			break
		}
		h.registrations[registration.InviteeID] = registration
		h.inviteIDs = append(h.inviteIDs, registration.InviteeID)
	}
	return nil
}

func (h *queryHandler) getSuperiores(ctx context.Context) error {
	if h.Limit < int32(len(h.inviteIDs)) {
		return fmt.Errorf("limit should be greater than userids")
	}

	resetInviteIDs := false
	if h.Offset == 0 {
		h.Limit -= int32(len(h.inviteIDs))
	} else if h.Offset < int32(len(h.inviteIDs)) {
		return fmt.Errorf("offset should be greater than userids")
	} else {
		h.Offset -= int32(len(h.inviteIDs))
		resetInviteIDs = true
	}

	registrations, total, err := registrationmwcli.GetSuperiores(ctx, &registrationmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		InviteeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
	}, h.Offset, h.Limit)
	if err != nil {
		return err
	}

	h.total = total + uint32(len(h.inviteIDs))
	if resetInviteIDs {
		h.inviteIDs = []string{}
	}
	for i, registration := range registrations {
		if int32(i) >= h.Limit {
			break
		}
		h.registrations[registration.InviteeID] = registration
		h.inviteIDs = append(h.inviteIDs, registration.InviterID)
	}

	offset := int32(0)
	limit := constant.DefaultRowLimit
	for {
		registrations, _, err = registrationmwcli.GetRegistrations(ctx, &registrationmwpb.Conds{
			AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			InviterIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(registrations) == 0 {
			break
		}
		for _, registration := range registrations {
			h.registrations[registration.InviteeID] = registration
			h.inviteIDs = append(h.inviteIDs, registration.InviterID)
		}
		offset += limit
	}
	return nil
}

func (h *queryHandler) getRegistrations(ctx context.Context) error {
	if h.AppID == nil {
		return fmt.Errorf("invalid appid")
	}
	if h.UserID != nil {
		h.inviteIDs = append(h.inviteIDs, *h.UserID)
		return h.getInvitees(ctx)
	} else if h.UserIDs != nil {
		h.inviteIDs = append(h.inviteIDs, *h.UserIDs...)
		return h.getSuperiores(ctx)
	}
	return fmt.Errorf("invalid userid")
}

func (h *queryHandler) getInviteesCount(ctx context.Context) error {
	if len(h.inviteIDs) == 0 {
		return nil
	}

	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		registrations, _, err := registrationmwcli.GetRegistrations(ctx, &registrationmwpb.Conds{
			AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			InviterIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(registrations) == 0 {
			break
		}
		for _, registration := range registrations {
			h.inviteesCount[registration.InviterID] += 1
		}
		offset += limit
	}

	return nil
}

func (h *queryHandler) getAchievements(ctx context.Context) error {
	if len(h.inviteIDs) == 0 {
		return nil
	}

	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		achievements, _, err := achievementmwcli.GetAchievements(ctx, &achievementmwpb.Conds{
			AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(achievements) == 0 {
			break
		}
		for _, achievement := range achievements {
			h.achievements[achievement.UserID] = achievement
		}
		offset += limit
	}
	return nil
}

func (h *queryHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, achievement := range h.achievements {
		coinTypeIDs = append(coinTypeIDs, achievement.CoinTypeID)
	}
	for _, good := range h.goods {
		coinTypeIDs = append(coinTypeIDs, good.CoinTypeID)
	}
	coins, _, err := appcoinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		AppID:       &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: coinTypeIDs},
	}, 0, int32(len(coinTypeIDs)))
	if err != nil {
		return err
	}
	for _, coin := range coins {
		h.coins[coin.CoinTypeID] = coin
	}
	return nil
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	if len(h.inviteIDs) == 0 {
		return nil
	}
	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		IDs:   &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
	}, 0, int32(len(h.inviteIDs)))
	if err != nil {
		return err
	}
	for _, user := range users {
		h.users[user.ID] = user
	}
	return nil
}

func (h *queryHandler) getGoods(ctx context.Context) error {
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(goods) == 0 {
			break
		}
		for _, good := range goods {
			if !good.Visible {
				continue
			}
			h.goods[good.GoodID] = good
		}
		offset += limit
	}
	return nil
}

func (h *queryHandler) getCommissions(ctx context.Context) error {
	if len(h.inviteIDs) == 0 {
		return nil
	}

	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		commissions, _, err := commmwcli.GetCommissions(ctx, &commissionmwpb.Conds{
			AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserIDs:    &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
			SettleType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(types.SettleType_GoodOrderPayment)},
			EndAt:      &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(commissions) == 0 {
			break
		}
		for _, commission := range commissions {
			commissions, ok := h.commissions[commission.GoodID]
			if !ok {
				commissions = map[string]*commissionmwpb.Commission{}
			}
			commissions[commission.UserID] = commission
			h.commissions[commission.GoodID] = commissions
		}
		offset += limit
	}
	return nil
}

func (h *queryHandler) formalizeUsers() {
	for _, user := range h.users {
		invitedAt := uint32(0)
		var inviterID *string

		registration, ok := h.registrations[user.ID]
		if ok {
			invitedAt = registration.CreatedAt
			inviterID = &registration.InviterID
		}

		h.infoMap[user.ID] = &npool.Achievement{
			InviterID:     inviterID,
			UserID:        user.ID,
			Username:      user.Username,
			EmailAddress:  user.EmailAddress,
			PhoneNO:       user.PhoneNO,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Kol:           user.Kol,
			TotalInvitees: h.inviteesCount[user.ID],
			CreatedAt:     user.CreatedAt,
			InvitedAt:     invitedAt,
		}
	}
}

func (h *queryHandler) userGoodCommission(appID, goodID, userID string) *commissionmwpb.Commission {
	commissions, ok := h.commissions[goodID]
	if ok {
		commission, ok := commissions[userID]
		if ok {
			return commission
		}
	}
	return &commissionmwpb.Commission{
		AppID:            appID,
		UserID:           userID,
		GoodID:           goodID,
		AmountOrPercent:  decimal.NewFromInt(0).String(),
		SettleType:       types.SettleType_GoodOrderPayment,
		SettleMode:       types.SettleMode_DefaultSettleMode,
		SettleAmountType: types.SettleAmountType_DefaultSettleAmountType,
		SettleInterval:   types.SettleInterval_DefaultSettleInterval,
	}
}

func (h *queryHandler) formalizeAchievements() {
	for _, achievement := range h.achievements {
		info, ok := h.infoMap[achievement.UserID]
		if !ok {
			continue
		}
		coin, ok := h.coins[achievement.CoinTypeID]
		if !ok {
			continue
		}
		good, ok := h.goods[achievement.GoodID]
		if !ok {
			continue
		}
		if good.CoinTypeID != achievement.CoinTypeID {
			logger.Sugar().Warnw(
				"formalizeAchievements",
				"AchievementID", achievement.ID,
				"GoodID", achievement.GoodID,
				"GoodCoinTypeID", good.CoinTypeID,
				"CoinTypeID", achievement.CoinTypeID,
				"State", "Achievement cointypeid is not good cointypeid",
			)
			continue
		}

		commission := h.userGoodCommission(achievement.AppID, achievement.GoodID, achievement.UserID)

		info.Achievements = append(info.Achievements, &npool.GoodAchievement{
			GoodID:                     achievement.GoodID,
			GoodName:                   good.GoodName,
			GoodUnit:                   good.Unit,
			CommissionValue:            commission.AmountOrPercent,
			CommissionThreshold:        commission.Threshold,
			CommissionSettleType:       commission.SettleType,
			CommissionSettleMode:       commission.SettleMode,
			CommissionSettleAmountType: commission.SettleAmountType,
			CommissionSettleInterval:   commission.SettleInterval,
			CoinTypeID:                 coin.CoinTypeID,
			CoinName:                   coin.Name,
			CoinLogo:                   coin.Logo,
			CoinUnit:                   coin.Unit,
			TotalUnits:                 achievement.TotalUnits,
			SelfUnits:                  achievement.SelfUnits,
			TotalAmount:                achievement.TotalAmount,
			SelfAmount:                 achievement.SelfAmount,
			TotalCommission:            achievement.TotalCommission,
			SelfCommission:             achievement.SelfCommission,
		})
		h.infoMap[achievement.UserID] = info
		achievedGoods, ok := h.achievedGoods[achievement.GoodID]
		if !ok {
			achievedGoods = map[string]struct{}{}
		}
		achievedGoods[achievement.UserID] = struct{}{}
		h.achievedGoods[achievement.GoodID] = achievedGoods
	}
}

func (h *queryHandler) formalizeNew() {
	for _, user := range h.users {
		for _, good := range h.goods {
			achievedGoods, ok := h.achievedGoods[good.GoodID]
			if ok {
				if _, ok := achievedGoods[user.ID]; ok {
					continue
				}
			}
			coin, ok := h.coins[good.CoinTypeID]
			if !ok {
				logger.Sugar().Warnw(
					"formalizeNew",
					"UserID", user.ID,
					"CoinTypeID", good.CoinTypeID,
					"State", "Invalid coin",
				)
				continue
			}
			info, ok := h.infoMap[user.ID]
			if !ok {
				logger.Sugar().Warnw(
					"formalizeNew",
					"UserID", user.ID,
					"State", "We should have info here",
				)
				continue
			}

			commission := h.userGoodCommission(good.AppID, good.GoodID, user.ID)

			info.Achievements = append(info.Achievements, &npool.GoodAchievement{
				GoodID:                     good.GoodID,
				GoodName:                   good.GoodName,
				GoodUnit:                   good.Unit,
				CommissionValue:            commission.AmountOrPercent,
				CommissionThreshold:        commission.Threshold,
				CommissionSettleType:       commission.SettleType,
				CommissionSettleMode:       commission.SettleMode,
				CommissionSettleAmountType: commission.SettleAmountType,
				CommissionSettleInterval:   commission.SettleInterval,
				CoinTypeID:                 coin.CoinTypeID,
				CoinName:                   coin.Name,
				CoinLogo:                   coin.Logo,
				CoinUnit:                   coin.Unit,
				TotalAmount:                decimal.NewFromInt(0).String(),
				SelfAmount:                 decimal.NewFromInt(0).String(),
				TotalUnits:                 decimal.NewFromInt(0).String(),
				SelfUnits:                  decimal.NewFromInt(0).String(),
				TotalCommission:            decimal.NewFromInt(0).String(),
				SelfCommission:             decimal.NewFromInt(0).String(),
			})
			h.infoMap[user.ID] = info
		}
	}
}

func (h *queryHandler) getStatements(ctx context.Context) error {
	if h.UserID == nil {
		return nil
	}
	offset := int32(0)
	limit := constant.DefaultRowLimit
	for {
		statements, _, err := statementmwcli.GetStatements(ctx, &statementmwpb.Conds{
			AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(statements) == 0 {
			break
		}
		h.statements = append(h.statements, statements...)
		offset += limit
	}
	return nil
}

func (h *queryHandler) formalizeDirectContribution(ctx context.Context) error {
	if err := h.getStatements(ctx); err != nil {
		return err
	}
	for _, statement := range h.statements {
		info, ok := h.infoMap[statement.DirectContributorID]
		if !ok {
			continue
		}
		for _, achievement := range info.Achievements {
			if achievement.GoodID != statement.GoodID {
				continue
			}
			amount, _ := decimal.NewFromString(statement.Commission)
			superior, _ := decimal.NewFromString(achievement.SuperiorCommission)
			currency, _ := decimal.NewFromString(statement.PaymentCoinUSDCurrency)
			achievement.SuperiorCommission = superior.Add(amount.Mul(currency)).String()
			break
		}
	}
	return nil
}

func (h *queryHandler) formalize(ctx context.Context) error {
	h.formalizeUsers()
	h.formalizeAchievements()
	h.formalizeNew()
	if err := h.formalizeDirectContribution(ctx); err != nil {
		return err
	}
	for _, info := range h.infoMap {
		h.infos = append(h.infos, info)
	}
	return nil
}

func (h *Handler) GetAchievements(ctx context.Context) ([]*npool.Achievement, uint32, error) {
	handler := &queryHandler{
		Handler:       h,
		registrations: map[string]*registrationmwpb.Registration{},
		inviteIDs:     []string{},
		achievements:  map[string]*achievementmwpb.Achievement{},
		inviteesCount: map[string]uint32{},
		coins:         map[string]*appcoinmwpb.Coin{},
		users:         map[string]*usermwpb.User{},
		goods:         map[string]*appgoodmwpb.Good{},
		commissions:   map[string]map[string]*commissionmwpb.Commission{},
		achievedGoods: map[string]map[string]struct{}{},
		statements:    []*statementmwpb.Statement{},
		infoMap:       map[string]*npool.Achievement{},
		infos:         []*npool.Achievement{},
	}
	if err := handler.getRegistrations(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getInviteesCount(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getAchievements(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCommissions(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.formalize(ctx); err != nil {
		return nil, 0, err
	}
	return handler.infos, handler.total, nil
}
