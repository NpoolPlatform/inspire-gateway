//nolint:dupl
package commission

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/appgood"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
)

type queryHandler struct {
	*Handler
	users map[string]*usermwpb.User
	goods map[string]*appgoodmwpb.Good
	coins map[string]*coinmwpb.Coin
	comms []*commmwpb.Commission
	infos []*npool.Commission
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	userIDs := []string{}
	for _, comm := range h.comms {
		userIDs = append(userIDs, comm.UserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.comms[0].AppID},
		IDs:   &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, int32(0), int32(len(userIDs)))
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
	for _, comm := range h.comms {
		goodIDs = append(goodIDs, comm.GoodID)
	}
	if len(goodIDs) == 0 {
		return nil
	}

	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmgrpb.Conds{
		AppID:   &commonpb.StringVal{Op: cruder.EQ, Value: h.comms[0].AppID},
		GoodIDs: &commonpb.StringSliceVal{Op: cruder.IN, Value: goodIDs},
	}, int32(0), int32(len(goodIDs)))
	if err != nil {
		return err
	}

	for _, good := range goods {
		h.goods[good.GoodID] = good
	}
	return nil
}

func (h *queryHandler) getCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, good := range h.goods {
		coinTypeIDs = append(coinTypeIDs, good.CoinTypeID)
	}
	if len(coinTypeIDs) == 0 {
		return nil
	}

	coins, _, err := coinmwcli.GetCoins(ctx, &coinmwpb.Conds{
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

		info := &npool.Commission{
			ID:              comm.ID,
			AppID:           comm.AppID,
			UserID:          comm.UserID,
			Username:        user.Username,
			EmailAddress:    user.EmailAddress,
			PhoneNO:         user.PhoneNO,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Kol:             user.Kol,
			SettleType:      comm.SettleType,
			SettleMode:      comm.SettleMode,
			SettleInterval:  comm.SettleInterval,
			GoodID:          comm.GoodID,
			AmountOrPercent: comm.AmountOrPercent,
			Threshold:       comm.Threshold,
			StartAt:         comm.StartAt,
			EndAt:           comm.EndAt,
			CreatedAt:       comm.CreatedAt,
			UpdatedAt:       comm.UpdatedAt,
		}

		good, ok := h.goods[comm.GoodID]
		if !ok {
			continue
		}
		coin, ok := h.coins[good.CoinTypeID]
		if !ok {
			continue
		}

		info.GoodName = good.GoodName
		info.CoinTypeID = good.CoinTypeID
		info.CoinName = coin.Name
		info.CoinLogo = coin.Logo
		h.infos = append(h.infos, info)
	}
}

func (h *Handler) GetCommission(ctx context.Context) (*npool.Commission, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := commmwcli.GetCommission(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler: h,
		users:   map[string]*usermwpb.User{},
		goods:   map[string]*appgoodmwpb.Good{},
		coins:   map[string]*coinmwpb.Coin{},
		comms:   []*commmwpb.Commission{info},
		infos:   []*npool.Commission{},
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, err
	}
	if err := handler.getGoods(ctx); err != nil {
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

func (h *Handler) GetCommissions(ctx context.Context) ([]*npool.Commission, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appid")
	}
	conds := &commmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
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

	handler := &queryHandler{
		Handler: h,
		users:   map[string]*usermwpb.User{},
		goods:   map[string]*appgoodmwpb.Good{},
		coins:   map[string]*coinmwpb.Coin{},
		comms:   infos,
		infos:   []*npool.Commission{},
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
	handler.formalize()
	return handler.infos, total, nil
}
