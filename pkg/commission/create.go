package commission

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	commmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/appcoin"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	"github.com/shopspring/decimal"
)

func CreateCommission(
	ctx context.Context,
	appID, userID string,
	goodID *string,
	settleType commmgrpb.SettleType,
	value decimal.Decimal,
	startAt *uint32,
) (
	*npool.Commission,
	error,
) {
	user, err := usermwcli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	if goodID != nil {
		good, err := goodmwcli.GetGoodOnly(ctx, &goodmgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
			GoodID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: *goodID,
			},
		})
		if err != nil {
			return nil, err
		}
		if good == nil {
			return nil, fmt.Errorf("invalid good")
		}

		coin, err := appcoinmwcli.GetCoinOnly(ctx, &appcoinmwpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: appID,
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
	}

	req := &commmwpb.CommissionReq{
		AppID:      &appID,
		UserID:     &userID,
		GoodID:     goodID,
		SettleType: &settleType,
		StartAt:    startAt,
	}

	valueStr := value.String()

	switch settleType {
	case commmgrpb.SettleType_GoodOrderPercent:
		req.Percent = &valueStr
	case commmgrpb.SettleType_LimitedOrderPercent:
		fallthrough //nolint
	case commmgrpb.SettleType_AmountThreshold:
		fallthrough //nolint
	case commmgrpb.SettleType_NoCommission:
		return nil, fmt.Errorf("not implemented")
	default:
		return nil, fmt.Errorf("unknown settle type")
	}

	info, err := commmwcli.CreateCommission(ctx, req)
	if err != nil {
		return nil, err
	}

	return GetCommission(ctx, info.ID, settleType)
}
