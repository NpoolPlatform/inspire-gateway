package commission

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	commmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	sendmwpb "github.com/NpoolPlatform/message/npool/third/mw/v1/send"
	sendmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/send"

	tmplmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	tmplmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	chanmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	"github.com/shopspring/decimal"
)

//nolint
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
	if !user.Kol {
		return nil, fmt.Errorf("user not kol")
	}

	if goodID != nil {
		good, err := goodmwcli.GetGoodOnly(ctx, &goodmgrpb.Conds{
			AppID:  &commonpb.StringVal{Op: cruder.EQ, Value: appID},
			GoodID: &commonpb.StringVal{Op: cruder.EQ, Value: *goodID},
		})
		if err != nil {
			return nil, err
		}
		if good == nil {
			return nil, fmt.Errorf("invalid good")
		}

		coin, err := appcoinmwcli.GetCoinOnly(ctx, &appcoinmwpb.Conds{
			AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: appID},
			CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: good.CoinTypeID},
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
		fallthrough //nolint
	case commmgrpb.SettleType_GoodOrderValuePercent:
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

	comm, err := GetCommission(ctx, info.ID, settleType)
	if err != nil {
		return nil, err
	}

	lang, err := applangmwcli.GetLangOnly(ctx, &applangmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: appID},
		Main:  &basetypes.BoolVal{Op: cruder.EQ, Value: true},
	})
	if err != nil {
		logger.Sugar().Warnw("CreateCommission", "Error", err)
		return comm, nil
	}
	if lang == nil {
		logger.Sugar().Warnw("CreateCommission", "Error", "Main AppLang not exist")
		return comm, nil
	}

	info1, err := tmplmwcli.GenerateText(ctx, &tmplmwpb.GenerateTextRequest{
		AppID:     appID,
		LangID:    lang.LangID,
		Channel:   chanmgrpb.NotifChannel_ChannelEmail,
		EventType: basetypes.UsedFor_SetCommission,
	})
	if err != nil {
		logger.Sugar().Warnw("CreateCommission", "Error", err)
		return comm, nil
	}
	if info1 == nil {
		logger.Sugar().Warnw("CreateCommission", "Error", "Cannot generate text")
		return comm, nil
	}

	err = sendmwcli.SendMessage(ctx, &sendmwpb.SendMessageRequest{
		Subject:     info1.Subject,
		Content:     info1.Content,
		From:        info1.From,
		To:          user.EmailAddress,
		ToCCs:       info1.ToCCs,
		ReplyTos:    info1.ReplyTos,
		AccountType: basetypes.SignMethod_Email,
	})
	if err != nil {
		logger.Sugar().Warnw("CreateCommission", "Error", "Cannot send message")
		return comm, nil
	}

	return comm, nil
}
