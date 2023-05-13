//nolint:dupl
package commission

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/appgood"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/appcoin"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func GetCommission(ctx context.Context, id string, settleType mgrpb.SettleType) (*npool.Commission, error) {
	info, err := commmwcli.GetCommission(ctx, id, settleType)
	if err != nil {
		return nil, err
	}

	user, err := usermwcli.GetUser(ctx, info.AppID, info.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	comm := &npool.Commission{
		ID:             info.ID,
		AppID:          info.AppID,
		UserID:         info.UserID,
		Username:       user.Username,
		EmailAddress:   user.EmailAddress,
		PhoneNO:        user.PhoneNO,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Kol:            user.Kol,
		SettleType:     info.SettleType,
		SettleMode:     info.SettleMode,
		SettleInterval: info.SettleInterval,
		GoodID:         info.GoodID,
		Percent:        info.Percent,
		Amount:         info.Amount,
		Threshold:      info.Threshold,
		StartAt:        info.StartAt,
		EndAt:          info.EndAt,
		CreatedAt:      info.CreatedAt,
		UpdatedAt:      info.UpdatedAt,
	}

	if info.GoodID != nil {
		good, err := goodmwcli.GetGoodOnly(ctx, &goodmgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: info.AppID,
			},
			GoodID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: info.GetGoodID(),
			},
		})
		if err != nil {
			return nil, err
		}
		if good == nil {
			return nil, fmt.Errorf("invalid good")
		}

		coin, err := coinmwcli.GetCoinOnly(ctx, &coinmwpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: info.AppID,
			},
			CoinTypeID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: good.CoinTypeID,
			},
		})
		if err != nil {
			return nil, err
		}
		if coin == nil {
			return nil, fmt.Errorf("invalid coin")
		}

		comm.GoodName = &good.GoodName
		comm.CoinTypeID = &good.CoinTypeID
		comm.CoinName = &coin.Name
		comm.CoinLogo = &coin.Logo
	}

	return comm, nil
}

func GetCommissions(ctx context.Context, conds *commmwpb.Conds, offset, limit int32) ([]*npool.Commission, uint32, error) { //nolint
	infos, total, err := commmwcli.GetCommissions(ctx, conds, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	userIDs := []string{}
	for _, info := range infos {
		userIDs = append(userIDs, info.UserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, 0, int32(len(userIDs)))
	if err != nil {
		return nil, 0, err
	}

	userMap := map[string]*usermwpb.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	goodIDs := []string{}
	for _, info := range infos {
		if info.GoodID != nil {
			goodIDs = append(goodIDs, info.GetGoodID())
		}
	}

	goods, _, err := goodmwcli.GetGoods(ctx, &goodmgrpb.Conds{
		AppID: conds.AppID,
		GoodIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: goodIDs,
		},
	}, int32(0), int32(len(goodIDs)))
	if err != nil {
		return nil, 0, err
	}

	goodMap := map[string]*goodmwpb.Good{}
	for _, good := range goods {
		goodMap[good.GoodID] = good
	}

	coinTypeIDs := []string{}
	for _, good := range goods {
		coinTypeIDs = append(coinTypeIDs, good.CoinTypeID)
	}

	coins, _, err := coinmwcli.GetCoins(ctx, &coinmwpb.Conds{
		AppID: conds.AppID,
		CoinTypeIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: coinTypeIDs,
		},
	}, int32(0), int32(len(coinTypeIDs)))
	if err != nil {
		return nil, 0, err
	}

	coinMap := map[string]*coinmwpb.Coin{}
	for _, coin := range coins {
		coinMap[coin.CoinTypeID] = coin
	}

	comms := []*npool.Commission{}
	for _, info := range infos {
		user, ok := userMap[info.UserID]
		if !ok {
			continue
		}

		comm := &npool.Commission{
			ID:             info.ID,
			AppID:          info.AppID,
			UserID:         info.UserID,
			Username:       user.Username,
			EmailAddress:   user.EmailAddress,
			PhoneNO:        user.PhoneNO,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Kol:            user.Kol,
			SettleType:     info.SettleType,
			SettleMode:     info.SettleMode,
			SettleInterval: info.SettleInterval,
			GoodID:         info.GoodID,
			Percent:        info.Percent,
			Amount:         info.Amount,
			Threshold:      info.Threshold,
			StartAt:        info.StartAt,
			EndAt:          info.EndAt,
			CreatedAt:      info.CreatedAt,
			UpdatedAt:      info.UpdatedAt,
		}

		if info.GoodID == nil {
			comms = append(comms, comm)
			continue
		}

		good, ok := goodMap[info.GetGoodID()]
		if !ok {
			continue
		}

		coin, ok := coinMap[good.CoinTypeID]
		if !ok {
			continue
		}

		comm.GoodName = &good.GoodName
		comm.CoinTypeID = &good.CoinTypeID
		comm.CoinName = &coin.Name
		comm.CoinLogo = &coin.Logo

		comms = append(comms, comm)
	}

	return comms, total, nil
}
