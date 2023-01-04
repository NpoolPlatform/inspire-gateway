package commission

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/appcoin"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
)

func GetCommission(ctx context.Context, id string) (*npool.Commission, error) {
	info, err := commmwcli.GetCommission(ctx, id)
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
