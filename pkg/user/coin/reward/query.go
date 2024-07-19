package reward

import (
	"context"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	usercoinrewardmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/user/coin/reward"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/user/coin/reward"
	usercoinrewardmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/user/coin/reward"
)

type queryHandler struct {
	*Handler
	appcoin         map[string]*appcoinmwpb.Coin
	infos           []*npool.UserCoinReward
	usercoinrewards []*usercoinrewardmwpb.UserCoinReward
}

func (h *queryHandler) getAppCoins(ctx context.Context) error {
	coinTypeIDs := []string{}
	for _, val := range h.usercoinrewards {
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
	for _, reward := range h.usercoinrewards {
		coin, ok := h.appcoin[reward.CoinTypeID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.UserCoinReward{
			ID:           reward.ID,
			EntID:        reward.EntID,
			AppID:        reward.AppID,
			UserID:       reward.UserID,
			CoinTypeID:   reward.CoinTypeID,
			CoinRewards:  reward.CoinRewards,
			CoinName:     coin.CoinName,
			DisplayNames: coin.DisplayNames,
			CoinLogo:     coin.Logo,
			CoinUnit:     coin.Unit,
			CreatedAt:    reward.CreatedAt,
			UpdatedAt:    reward.UpdatedAt,
		})
	}
}

func (h *Handler) GetUserCoinReward(ctx context.Context) (*npool.UserCoinReward, error) {
	if h.EntID == nil {
		return nil, wlog.Errorf("invalid entid")
	}

	info, err := usercoinrewardmwcli.GetUserCoinReward(ctx, *h.EntID)
	if err != nil {
		return nil, wlog.WrapError(err)
	}
	if info == nil {
		return nil, nil
	}
	handler := &queryHandler{
		Handler: h,
		appcoin: map[string]*appcoinmwpb.Coin{},
		infos:   []*npool.UserCoinReward{},
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

func (h *Handler) GetUserCoinRewards(ctx context.Context) ([]*npool.UserCoinReward, uint32, error) {
	conds := &usercoinrewardmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}
	infos, _, err := usercoinrewardmwcli.GetUserCoinRewards(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, wlog.WrapError(err)
	}
	handler := &queryHandler{
		Handler:         h,
		appcoin:         map[string]*appcoinmwpb.Coin{},
		usercoinrewards: []*usercoinrewardmwpb.UserCoinReward{},
		infos:           []*npool.UserCoinReward{},
	}
	handler.usercoinrewards = infos

	if err := handler.getAppCoins(ctx); err != nil {
		return nil, 0, wlog.WrapError(err)
	}

	handler.formalize()

	return handler.infos, uint32(len(handler.infos)), nil
}
