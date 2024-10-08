package achievement

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	goodcoinmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good/coin"
	requiredgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good/required"
	powerrentalmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/powerrental"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	goodachievementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement/good"
	coinachievementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement/good/coin"
	orderstatementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement/statement/order"
	achievementusermwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement/user"
	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	registrationmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	goodtypes "github.com/NpoolPlatform/message/npool/basetypes/good/v1"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	goodcoinmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good/coin"
	requiredgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good/required"
	powerrentalmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/powerrental"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/achievement"
	goodachievementmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement/good"
	coinachievementmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement/good/coin"
	orderstatementmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement/statement/order"
	achievementusermwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement/user"
	commissionmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
	registrationmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type queryHandler struct {
	*Handler
	registrations     map[string]*registrationmwpb.Registration
	inviteIDs         []string
	coinAchievements  map[string]map[string]*coinachievementmwpb.Achievement // userid->goodcointypeid->achievement
	goodAchievements  []*goodachievementmwpb.Achievement
	inviteesCount     map[string]uint32
	coins             map[string]*appcoinmwpb.Coin
	users             map[string]*usermwpb.User
	appGoods          map[string]*appgoodmwpb.Good
	goodQuantityUnits map[string]string
	requiredGoods     map[string]*requiredgoodmwpb.Required
	goodMainCoins     map[string]*goodcoinmwpb.GoodCoin
	commissions       map[string]map[string]*commissionmwpb.Commission
	total             uint32
	achievedGoods     map[string]map[string]struct{}
	achievementUsers  map[string]*achievementusermwpb.AchievementUser
	statements        []*orderstatementmwpb.Statement
	infoMap           map[string]*npool.Achievement
	infos             []*npool.Achievement
}

func (h *Handler) checkUser(ctx context.Context) error {
	if h.UserID != nil {
		user, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("invalid user id")
		}
	}
	return nil
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

func (h *queryHandler) getCoinAchievements(ctx context.Context) error {
	if len(h.inviteIDs) == 0 {
		return nil
	}

	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		achievements, _, err := coinachievementmwcli.GetAchievements(ctx, &coinachievementmwpb.Conds{
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
			coinAchievementMap := map[string]*coinachievementmwpb.Achievement{
				achievement.GoodCoinTypeID: achievement,
			}
			h.coinAchievements[achievement.UserID] = coinAchievementMap
		}
		offset += limit
	}
	return nil
}

func (h *queryHandler) getGoodAchievements(ctx context.Context) error {
	if len(h.inviteIDs) == 0 {
		return nil
	}

	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		achievements, _, err := goodachievementmwcli.GetAchievements(ctx, &goodachievementmwpb.Conds{
			AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(achievements) == 0 {
			break
		}
		h.goodAchievements = append(h.goodAchievements, achievements...)
		offset += limit
	}
	return nil
}

func (h *queryHandler) getAchievementUsers(ctx context.Context) error {
	if len(h.inviteIDs) == 0 {
		return nil
	}

	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		achievements, _, err := achievementusermwcli.GetAchievementUsers(ctx, &achievementusermwpb.Conds{
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
			h.achievementUsers[achievement.UserID] = achievement
		}
		offset += limit
	}
	return nil
}

func (h *queryHandler) getGoodCoins(ctx context.Context) (err error) {
	h.goodMainCoins = map[string]*goodcoinmwpb.GoodCoin{}
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		goodCoins, _, err := goodcoinmwcli.GetGoodCoins(ctx, &goodcoinmwpb.Conds{
			GoodIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: func() (_goodIDs []string) {
				for _, appGood := range h.appGoods {
					_goodIDs = append(_goodIDs, appGood.GoodID)
				}
				for _, required := range h.requiredGoods {
					_goodIDs = append(_goodIDs, required.MainGoodID, required.RequiredGoodID)
				}
				return
			}()},
			Main: &basetypes.BoolVal{Op: cruder.EQ, Value: true},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(goodCoins) == 0 {
			return nil
		}
		offset += limit
		for _, goodCoin := range goodCoins {
			h.goodMainCoins[goodCoin.GoodID] = goodCoin
		}
	}
}

func (h *queryHandler) getRequiredGoods(ctx context.Context) (err error) {
	h.requiredGoods = map[string]*requiredgoodmwpb.Required{}
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		requireds, _, err := requiredgoodmwcli.GetRequireds(ctx, &requiredgoodmwpb.Conds{
			GoodIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: func() (_goodIDs []string) {
				for _, appGood := range h.appGoods {
					_goodIDs = append(_goodIDs, appGood.GoodID)
				}
				return
			}()},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(requireds) == 0 {
			return nil
		}
		for _, required := range requireds {
			h.requiredGoods[required.RequiredGoodID] = required
		}
		offset += limit
	}
}

func (h *queryHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, achievementMap := range h.coinAchievements {
		for goodCoinTypeID := range achievementMap {
			if _, err := uuid.Parse(goodCoinTypeID); err != nil {
				continue
			}
			coinTypeIDs = append(coinTypeIDs, goodCoinTypeID)
		}
	}
	for _, goodCoin := range h.goodMainCoins {
		if _, err := uuid.Parse(goodCoin.CoinTypeID); err != nil {
			continue
		}
		coinTypeIDs = append(coinTypeIDs, goodCoin.CoinTypeID)
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

func (h *queryHandler) getGoodQuantityUnits(ctx context.Context) error {
	powerRentalIDs := func() (_goodIDs []string) {
		for _, appGood := range h.appGoods {
			if appGood.GoodType != goodtypes.GoodType_PowerRental &&
				appGood.GoodType != goodtypes.GoodType_LegacyPowerRental {
				continue
			}
			_goodIDs = append(_goodIDs, appGood.GoodID)
		}
		return
	}()
	// TODO: other type should be added when it's implemented
	powerRentals, _, err := powerrentalmwcli.GetPowerRentals(ctx, &powerrentalmwpb.Conds{
		GoodIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: powerRentalIDs},
	}, 0, int32(len(powerRentalIDs)))
	if err != nil {
		return err
	}
	for _, powerRental := range powerRentals {
		h.goodQuantityUnits[powerRental.GoodID] = powerRental.QuantityUnit
	}
	return nil
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	if len(h.inviteIDs) == 0 {
		return nil
	}
	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
	}, 0, int32(len(h.inviteIDs)))
	if err != nil {
		return err
	}
	for _, user := range users {
		h.users[user.EntID] = user
	}
	return nil
}

func (h *queryHandler) getAppGoods(ctx context.Context) error {
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
			h.appGoods[good.EntID] = good
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
			commissions, ok := h.commissions[commission.AppGoodID]
			if !ok {
				commissions = map[string]*commissionmwpb.Commission{}
			}
			commissions[commission.UserID] = commission
			h.commissions[commission.AppGoodID] = commissions
		}
		offset += limit
	}
	return nil
}

func (h *queryHandler) formalizeUsers() {
	for _, user := range h.users {
		invitedAt := uint32(0)
		var inviterID *string

		registration, ok := h.registrations[user.EntID]
		if ok {
			invitedAt = registration.CreatedAt
			inviterID = &registration.InviterID
		}

		info := &npool.Achievement{
			InviterID:    inviterID,
			UserID:       user.EntID,
			Username:     user.Username,
			EmailAddress: user.EmailAddress,
			PhoneNO:      user.PhoneNO,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Kol:          user.Kol,
			CreatedAt:    user.CreatedAt,
			InvitedAt:    invitedAt,
		}

		achievementUser, ok := h.achievementUsers[user.EntID]
		if ok {
			info.TotalCommission = achievementUser.TotalCommission
			info.SelfCommission = achievementUser.SelfCommission
			info.DirectInvitees = achievementUser.DirectInvitees
			info.IndirectInvitees = achievementUser.IndirectInvitees
			info.DirectConsumeAmount = achievementUser.DirectConsumeAmount
			info.InviteeConsumeAmount = achievementUser.InviteeConsumeAmount
		}

		h.infoMap[user.EntID] = info
	}
}

func (h *queryHandler) userGoodCommission(appID, goodID, appGoodID, userID string) *commissionmwpb.Commission {
	commissions, ok := h.commissions[appGoodID]
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
		AppGoodID:        appGoodID,
		AmountOrPercent:  decimal.NewFromInt(0).String(),
		SettleType:       types.SettleType_GoodOrderPayment,
		SettleMode:       types.SettleMode_DefaultSettleMode,
		SettleAmountType: types.SettleAmountType_DefaultSettleAmountType,
		SettleInterval:   types.SettleInterval_DefaultSettleInterval,
	}
}

func (h *queryHandler) achievementMainGoodCoin(goodID string) (*goodcoinmwpb.GoodCoin, error) {
	if required, ok := h.requiredGoods[goodID]; ok {
		goodMainCoin, ok := h.goodMainCoins[required.MainGoodID]
		if !ok {
			return nil, fmt.Errorf("invalid goodmaincoin")
		}
		return goodMainCoin, nil
	}
	goodMainCoin, ok := h.goodMainCoins[goodID]
	if !ok {
		return nil, fmt.Errorf("invalid goodmaincoin")
	}
	return goodMainCoin, nil
}

func (h *queryHandler) formalizeAchievements() {
	for _, achievement := range h.goodAchievements {
		info, ok := h.infoMap[achievement.UserID]
		if !ok {
			continue
		}
		good, ok := h.appGoods[achievement.AppGoodID]
		if !ok {
			continue
		}
		goodMainCoin, err := h.achievementMainGoodCoin(achievement.GoodID)
		if err != nil {
			continue
		}
		coin, ok := h.coins[goodMainCoin.CoinTypeID]
		if !ok {
			continue
		}
		commission := h.userGoodCommission(
			achievement.AppID,
			achievement.GoodID,
			achievement.AppGoodID,
			achievement.UserID,
		)
		info.Achievements = append(info.Achievements, &npool.GoodAchievement{
			GoodID:                     achievement.GoodID,
			GoodName:                   good.GoodName,
			GoodUnit:                   h.goodQuantityUnits[good.GoodID],
			AppGoodID:                  good.EntID,
			AppGoodName:                good.AppGoodName,
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
			TotalAmount:                achievement.TotalAmountUSD,
			SelfAmount:                 achievement.SelfAmountUSD,
			TotalCommission:            achievement.TotalCommissionUSD,
			SelfCommission:             achievement.SelfCommissionUSD,
		})
		h.infoMap[achievement.UserID] = info
		achievedGoods, ok := h.achievedGoods[achievement.AppGoodID]
		if !ok {
			achievedGoods = map[string]struct{}{}
		}
		achievedGoods[achievement.UserID] = struct{}{}
		h.achievedGoods[achievement.AppGoodID] = achievedGoods
	}
}

func (h *queryHandler) formalizeNew() {
	for _, user := range h.users {
		for _, good := range h.appGoods {
			achievedGoods, ok := h.achievedGoods[good.EntID]
			if ok {
				if _, ok := achievedGoods[user.EntID]; ok {
					continue
				}
			}
			goodMainCoin, err := h.achievementMainGoodCoin(good.GoodID)
			if err != nil {
				continue
			}
			coin, ok := h.coins[goodMainCoin.CoinTypeID]
			if !ok {
				continue
			}
			info, ok := h.infoMap[user.EntID]
			if !ok {
				continue
			}

			commission := h.userGoodCommission(good.AppID, good.GoodID, good.EntID, user.EntID)
			info.Achievements = append(info.Achievements, &npool.GoodAchievement{
				GoodID:                     good.GoodID,
				GoodName:                   good.GoodName,
				GoodUnit:                   h.goodQuantityUnits[good.GoodID],
				AppGoodID:                  good.EntID,
				AppGoodName:                good.AppGoodName,
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
			h.infoMap[user.EntID] = info
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
		statements, _, err := orderstatementmwcli.GetStatements(ctx, &orderstatementmwpb.Conds{
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
			if achievement.AppGoodID != statement.AppGoodID {
				continue
			}
			amount, _ := decimal.NewFromString(statement.CommissionAmountUSD)
			superior, _ := decimal.NewFromString(achievement.SuperiorCommission)
			achievement.SuperiorCommission = superior.Add(amount).String()
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
	if err := h.checkUser(ctx); err != nil {
		return nil, 0, err
	}

	handler := &queryHandler{
		Handler:           h,
		registrations:     map[string]*registrationmwpb.Registration{},
		inviteIDs:         []string{},
		inviteesCount:     map[string]uint32{},
		coins:             map[string]*appcoinmwpb.Coin{},
		users:             map[string]*usermwpb.User{},
		appGoods:          map[string]*appgoodmwpb.Good{},
		goodQuantityUnits: map[string]string{},
		commissions:       map[string]map[string]*commissionmwpb.Commission{},
		achievedGoods:     map[string]map[string]struct{}{},
		statements:        []*orderstatementmwpb.Statement{},
		goodAchievements:  []*goodachievementmwpb.Achievement{},
		coinAchievements:  map[string]map[string]*coinachievementmwpb.Achievement{},
		infoMap:           map[string]*npool.Achievement{},
		infos:             []*npool.Achievement{},
		achievementUsers:  map[string]*achievementusermwpb.AchievementUser{},
	}
	if err := handler.getRegistrations(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getInviteesCount(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getGoodAchievements(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCoinAchievements(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getAchievementUsers(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getRequiredGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getGoodCoins(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCommissions(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getGoodQuantityUnits(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.formalize(ctx); err != nil {
		return nil, 0, err
	}
	return handler.infos, handler.total, nil
}
