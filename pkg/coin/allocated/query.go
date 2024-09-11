package allocated

import (
	"context"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	coinallocatedmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/allocated"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/allocated"
	coinallocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/allocated"
)

type queryHandler struct {
	*Handler
	coinallocateds []*coinallocatedmwpb.CoinAllocated
	appcoin        map[string]*appcoinmwpb.Coin
	appuser        map[string]*appusermwpb.User
	infos          []*npool.CoinAllocated
	total          uint32
}

func (h *queryHandler) getCoinAllocateds(ctx context.Context) error {
	conds := &coinallocatedmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}
	infos, total, err := coinallocatedmwcli.GetCoinAllocateds(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return wlog.WrapError(err)
	}
	h.total = total
	h.coinallocateds = infos
	return nil
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	userIDs := []string{}
	for _, allocated := range h.coinallocateds {
		userIDs = append(userIDs, allocated.UserID)
	}
	users, _, err := usermwcli.GetUsers(ctx, &appusermwpb.Conds{
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, 0, int32(len(userIDs)))
	if err != nil {
		return wlog.WrapError(err)
	}

	for _, user := range users {
		h.appuser[user.EntID] = user
	}
	return nil
}

func (h *queryHandler) getAppCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, val := range h.coinallocateds {
		coinTypeIDs = append(coinTypeIDs, val.CoinTypeID)
	}
	coins, _, err := appcoinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		AppID:       &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: coinTypeIDs},
	}, 0, int32(len(coinTypeIDs)))
	if err != nil {
		return wlog.WrapError(err)
	}

	for _, coin := range coins {
		h.appcoin[coin.CoinTypeID] = coin
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, info := range h.coinallocateds {
		coin, ok := h.appcoin[info.CoinTypeID]
		if !ok {
			continue
		}
		user, ok := h.appuser[info.UserID]
		if !ok {
			continue
		}

		h.infos = append(h.infos, &npool.CoinAllocated{
			ID:           info.ID,
			EntID:        info.EntID,
			AppID:        info.AppID,
			CoinConfigID: info.CoinConfigID,
			CoinTypeID:   coin.CoinTypeID,
			CoinName:     coin.CoinName,
			DisplayNames: coin.DisplayNames,
			CoinLogo:     coin.Logo,
			CoinUnit:     coin.Unit,
			CoinAmount:   info.Value,
			UserID:       user.EntID,
			PhoneNO:      user.PhoneNO,
			EmailAddress: user.EmailAddress,
			Extra:        info.Extra,
			CreatedAt:    info.CreatedAt,
			UpdatedAt:    info.UpdatedAt,
		})
	}
}

func (h *Handler) GetCoinAllocated(ctx context.Context) (*npool.CoinAllocated, error) {
	if h.EntID == nil {
		return nil, wlog.Errorf("invalid entid")
	}

	info, err := coinallocatedmwcli.GetCoinAllocated(ctx, *h.EntID)
	if err != nil {
		return nil, wlog.WrapError(err)
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:        h,
		coinallocateds: []*coinallocatedmwpb.CoinAllocated{info},
		appcoin:        map[string]*appcoinmwpb.Coin{},
		appuser:        map[string]*appusermwpb.User{},
	}

	if err := handler.getUsers(ctx); err != nil {
		return nil, wlog.WrapError(err)
	}

	if err := handler.getAppCoins(ctx); err != nil {
		return nil, wlog.WrapError(err)
	}

	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}

	return handler.infos[0], nil
}

func (h *Handler) GetCoinAllocateds(ctx context.Context) ([]*npool.CoinAllocated, uint32, error) {
	handler := &queryHandler{
		Handler:        h,
		coinallocateds: []*coinallocatedmwpb.CoinAllocated{},
		appcoin:        map[string]*appcoinmwpb.Coin{},
		appuser:        map[string]*appusermwpb.User{},
		total:          uint32(0),
	}

	if err := handler.getCoinAllocateds(ctx); err != nil {
		return nil, 0, wlog.WrapError(err)
	}

	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, wlog.WrapError(err)
	}

	if err := handler.getAppCoins(ctx); err != nil {
		return nil, 0, wlog.WrapError(err)
	}

	handler.formalize()
	return handler.infos, handler.total, nil
}
