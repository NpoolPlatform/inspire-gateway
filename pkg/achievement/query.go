//nolint:dupl
package achievement

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	"github.com/shopspring/decimal"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	achievementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement"
	statementmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/achievement/statement"
	achievementmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement"
	statementmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/achievement/statement"

	registrationmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	registrationmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/appgood"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/achievement"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
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
	commissions   map[string]map[string]*commmwpb.Commission
	total         uint32
	achievedGoods map[string]map[string]struct{}
	statements    []*statementmwpb.Statement
	infoMap       map[string]*npool.Achievement
	infos         []*npool.Achievement
}

func (h *queryHandler) getInvitees(ctx context.Context) error {
	registrations, _, err := registrationmwcli.GetRegistrations(ctx, &registrationmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		InviterIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
	}, h.Offset, h.Limit)
	if err != nil {
		return err
	}
	for _, registration := range registrations {
		h.registrations[registration.InviteeID] = registration
		h.inviteIDs = append(h.inviteIDs, registration.InviteeID)
	}
	return nil
}

func (h *queryHandler) getSuperiores(ctx context.Context) error {
	registrations, _, err := registrationmwcli.GetSuperiores(ctx, &registrationmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		InviteeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
	}, h.Offset, h.Limit)
	if err != nil {
		return err
	}
	for _, registration := range registrations {
		h.registrations[registration.InviteeID] = registration
		h.inviteIDs = append(h.inviteIDs, registration.InviteeID)
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
	achievements, total, err := achievementmwcli.GetAchievements(ctx, &achievementmwpb.Conds{
		AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
	}, 0, int32(len(h.inviteIDs)))
	if err != nil {
		return err
	}
	for _, achievement := range achievements {
		h.achievements[achievement.UserID] = achievement
	}
	h.total = total
	return nil
}

func (h *queryHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, achievement := range h.achievements {
		coinTypeIDs = append(coinTypeIDs, achievement.CoinTypeID)
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
	goodIDs := []string{}
	for _, achievement := range h.achievements {
		goodIDs = append(goodIDs, achievement.GoodID)
	}
	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmgrpb.Conds{
		AppID:   &commonpb.StringVal{Op: cruder.EQ, Value: *h.AppID},
		GoodIDs: &commonpb.StringSliceVal{Op: cruder.IN, Value: goodIDs},
	}, 0, int32(len(goodIDs)))
	if err != nil {
		return err
	}
	for _, good := range goods {
		h.goods[good.GoodID] = good
	}
	return nil
}

func (h *queryHandler) getCommissions(ctx context.Context) error {
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		commissions, _, err := commmwcli.GetCommissions(ctx, &commmwpb.Conds{
			AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: h.inviteIDs},
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
				commissions = map[string]*commmwpb.Commission{}
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
		var inviterID string

		registration, ok := h.registrations[user.ID]
		if ok {
			invitedAt = registration.CreatedAt
			inviterID = registration.InviterID
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
		percent := decimal.NewFromInt(0).String()
		commissions, ok := h.commissions[achievement.GoodID]
		if ok {
			commission, ok := commissions[achievement.UserID]
			if ok {
				percent = commission.AmountOrPercent
			}
		}
		info.Achievements = append(
			info.Achievements,
			&npool.GoodAchievement{
				GoodID:            achievement.GoodID,
				GoodName:          good.GoodName,
				GoodUnit:          good.Unit,
				CommissionPercent: percent,
				CoinTypeID:        coin.ID,
				CoinName:          coin.Name,
				CoinLogo:          coin.Logo,
				CoinUnit:          coin.Unit,
				TotalUnits:        achievement.TotalUnits,
				SelfUnits:         achievement.SelfUnits,
				TotalAmount:       achievement.TotalAmount,
				SelfAmount:        achievement.SelfAmount,
				TotalCommission:   achievement.TotalCommission,
				SelfCommission:    achievement.SelfCommission,
			},
		)
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
	for _, info := range h.infoMap {
	nextCommission:
		for goodID, commissions := range h.commissions {
			achievedGoods, ok := h.achievedGoods[goodID]
			if !ok {
				break
			}
			for userID := range commissions {
				if _, ok := achievedGoods[userID]; ok {
					break nextCommission
				}

				percent := decimal.NewFromInt(0).String()
				commissions, ok := h.commissions[goodID]
				if ok {
					commission, ok := commissions[userID]
					if ok {
						percent = commission.AmountOrPercent
					}
				}

				good, ok := h.goods[goodID]
				if !ok {
					continue
				}
				coin, ok := h.coins[good.CoinTypeID]
				if !ok {
					continue
				}

				info.Achievements = append(
					info.Achievements,
					&npool.GoodAchievement{
						GoodID:            goodID,
						GoodName:          good.GoodName,
						GoodUnit:          good.Unit,
						CommissionPercent: percent,
						CoinTypeID:        coin.ID,
						CoinName:          coin.Name,
						CoinLogo:          coin.Logo,
						CoinUnit:          coin.Unit,
						TotalAmount:       decimal.NewFromInt(0).String(),
						SelfAmount:        decimal.NewFromInt(0).String(),
						TotalUnits:        decimal.NewFromInt(0).String(),
						SelfUnits:         decimal.NewFromInt(0).String(),
						TotalCommission:   decimal.NewFromInt(0).String(),
						SelfCommission:    decimal.NewFromInt(0).String(),
					},
				)
				h.infoMap[userID] = info
			}
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
	handler := &queryHandler{
		Handler:       h,
		registrations: map[string]*registrationmwpb.Registration{},
		inviteIDs:     []string{},
		achievements:  map[string]*achievementmwpb.Achievement{},
		inviteesCount: map[string]uint32{},
		coins:         map[string]*appcoinmwpb.Coin{},
		users:         map[string]*usermwpb.User{},
		goods:         map[string]*appgoodmwpb.Good{},
		commissions:   map[string]map[string]*commmwpb.Commission{},
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
	if err := handler.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getGoods(ctx); err != nil {
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
