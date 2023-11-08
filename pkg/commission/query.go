package commission

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	registrationmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
	registrationmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"

	"github.com/google/uuid"
)

type queryHandler struct {
	*Handler
	users    map[string]*usermwpb.User
	invitees []*registrationmwpb.Registration
	appGoods map[string]*appgoodmwpb.Good
	goods    map[string]*goodmwpb.Good
	coins    map[string]*appcoinmwpb.Coin
	comms    []*commmwpb.Commission
	infos    []*npool.Commission
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	userIDs := []string{}
	for _, comm := range h.comms {
		userIDs = append(userIDs, comm.UserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: h.comms[0].AppID},
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, int32(0), int32(len(userIDs)))
	if err != nil {
		return err
	}

	for _, user := range users {
		h.users[user.EntID] = user
	}
	return nil
}

func (h *queryHandler) getAppGoods(ctx context.Context) error {
	goodIDs := []string{}
	for _, comm := range h.comms {
		if _, err := uuid.Parse(comm.AppGoodID); err != nil {
			continue
		}
		goodIDs = append(goodIDs, comm.AppGoodID)
	}
	if len(goodIDs) == 0 {
		return nil
	}

	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.comms[0].AppID},
		IDs:   &basetypes.StringSliceVal{Op: cruder.IN, Value: goodIDs},
	}, int32(0), int32(len(goodIDs)))
	if err != nil {
		return err
	}

	for _, good := range goods {
		h.appGoods[good.ID] = good
	}
	return nil
}

func (h *queryHandler) getGoods(ctx context.Context) error {
	ids := []string{}
	for _, comm := range h.comms {
		if _, err := uuid.Parse(comm.GoodID); err != nil {
			continue
		}
		ids = append(ids, comm.GoodID)
	}
	if len(ids) == 0 {
		return nil
	}

	goods, _, err := goodmwcli.GetGoods(ctx, &goodmwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, int32(0), int32(len(ids)))
	if err != nil {
		return err
	}

	for _, good := range goods {
		h.goods[good.ID] = good
	}
	return nil
}

func (h *queryHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, good := range h.appGoods {
		coinTypeIDs = append(coinTypeIDs, good.CoinTypeID)
	}
	if len(coinTypeIDs) == 0 {
		return nil
	}

	coins, _, err := coinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		AppID:       &basetypes.StringVal{Op: cruder.EQ, Value: h.comms[0].AppID},
		CoinTypeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: coinTypeIDs},
	}, int32(0), int32(len(coinTypeIDs)))
	if err != nil {
		return err
	}

	for _, coin := range coins {
		h.coins[coin.CoinTypeID] = coin
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, comm := range h.comms {
		user, ok := h.users[comm.UserID]
		if !ok {
			continue
		}
		good, ok := h.goods[comm.GoodID]
		if !ok {
			continue
		}
		appgood, ok := h.appGoods[comm.AppGoodID]
		if !ok {
			continue
		}
		coin, ok := h.coins[appgood.CoinTypeID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.Commission{
			ID:               comm.ID,
			EntID:            comm.EntID,
			AppID:            comm.AppID,
			UserID:           comm.UserID,
			Username:         user.Username,
			EmailAddress:     user.EmailAddress,
			PhoneNO:          user.PhoneNO,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Kol:              user.Kol,
			SettleType:       comm.SettleType,
			SettleMode:       comm.SettleMode,
			SettleAmountType: comm.SettleAmountType,
			SettleInterval:   comm.SettleInterval,
			GoodID:           comm.GoodID,
			GoodTitle:        good.Title,
			AppGoodID:        comm.AppGoodID,
			GoodName:         appgood.GoodName,
			AmountOrPercent:  comm.AmountOrPercent,
			Threshold:        comm.Threshold,
			StartAt:          comm.StartAt,
			EndAt:            comm.EndAt,
			CoinTypeID:       appgood.CoinTypeID,
			CoinName:         coin.Name,
			CoinLogo:         coin.Logo,
			CreatedAt:        comm.CreatedAt,
			UpdatedAt:        comm.UpdatedAt,
		})
	}
}

func (h *Handler) GetCommission(ctx context.Context) (*npool.Commission, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}

	info, err := commmwcli.GetCommission(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:  h,
		users:    map[string]*usermwpb.User{},
		goods:    map[string]*goodmwpb.Good{},
		appGoods: map[string]*appgoodmwpb.Good{},
		coins:    map[string]*appcoinmwpb.Coin{},
		comms:    []*commmwpb.Commission{info},
		infos:    []*npool.Commission{},
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, err
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, err
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, err
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}

	return handler.infos[0], nil
}

func (h *queryHandler) getInvitees(ctx context.Context) error {
	if h.UserID == nil {
		return nil
	}

	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		infos, _, err := registrationmwcli.GetRegistrations(ctx, &registrationmwpb.Conds{
			AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			InviterID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(infos) == 0 {
			break
		}
		h.invitees = append(h.invitees, infos...)
		offset += limit
	}

	return nil
}

func (h *Handler) GetCommissions(ctx context.Context) ([]*npool.Commission, uint32, error) {
	handler := &queryHandler{
		Handler:  h,
		users:    map[string]*usermwpb.User{},
		goods:    map[string]*goodmwpb.Good{},
		appGoods: map[string]*appgoodmwpb.Good{},
		coins:    map[string]*appcoinmwpb.Coin{},
		infos:    []*npool.Commission{},
	}

	if err := handler.getInvitees(ctx); err != nil {
		return nil, 0, err
	}

	conds := &commmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		userIDs := []string{*h.UserID}
		for _, invitee := range handler.invitees {
			userIDs = append(userIDs, invitee.InviteeID)
		}
		conds.UserIDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs}
	}
	if h.EndAt != nil {
		conds.EndAt = &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.EndAt}
	}
	infos, total, err := commmwcli.GetCommissions(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}
	handler.comms = infos

	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getAppGoods(ctx); err != nil {
		return nil, 0, err
	}
	if err := handler.getCoins(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()
	return handler.infos, total, nil
}
