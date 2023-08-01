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

	"github.com/shopspring/decimal"
)

type createHandler struct {
	*Handler
	user *usermwpb.User
	req  *commmwpb.CommissionReq
}

func (h *createHandler) getUser(ctx context.Context) error {
	if h.AppID == nil {
		return fmt.Errorf("invalid appid")
	}
	if h.UserID == nil {
		return fmt.Errorf("invalid userid")
	}
	user, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("invalid user")
	}
	if !user.Kol {
		return fmt.Errorf("permission denied")
	}
	h.user = user
	return nil
}

func (h *createHandler) validateGood(ctx context.Context) error {
	if h.GoodID == nil {
		return nil
	}

	good, err := goodmwcli.GetGoodOnly(ctx, &goodmgrpb.Conds{
		AppID:  &commonpb.StringVal{Op: cruder.EQ, Value: *h.AppID},
		GoodID: &commonpb.StringVal{Op: cruder.EQ, Value: *h.GoodID},
	})
	if err != nil {
		return nil, err
	}
	if good == nil {
		return nil, fmt.Errorf("invalid good")
	}

	coin, err := appcoinmwcli.GetCoinOnly(ctx, &appcoinmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: good.CoinTypeID},
	})
	if err != nil {
		return nil, err
	}
	if coin == nil {
		return nil, fmt.Errorf("invalid coin")
	}

	return nil
}

func (h *createHandler) createCommission(ctx context.Context) error {
	h.req = &commmwpb.CommissionReq{
		AppID:           h.AppID,
		UserID:          h.UserID,
		GoodID:          h.GoodID,
		SettleType:      h.SettleType,
		SettleMode:      h.SettleMode,
		StartAt:         h.StartAt,
		AmountOrPercemt: h.AmountOrPercent,
	}
	info, err := commmwcli.CreateCommission(ctx, h.req)
	if err != nil {
		return err
	}
	h.ID = &info.ID
	return nil
}

func (h *createHandler) notifyCreateCommission(ctx context.Context) {
	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		return publisher.Update(
			basetypes.MsgID_CreateCommissionReq.String(),
			nil,
			nil,
			nil,
			&commmwpb.Commission{
				AppID:           *h.AppID,
				UserID:          *h.UserID,
				GoodID:          *h.GoodID,
				SettleType:      *h.SettleType,
				SettleMode:      *h.SettleMode,
				StartAt:         *h.StartAt,
				AmountOrPercemt: *h.AmountOrPercent,
			},
		)
	}); err != nil {
		logger.Sugar().Errorw(
			"rewardSignup",
			"AppID", h.AppID,
			"UserID", h.UserID,
			"Error", err,
		)
	}
}

func (h *Handler) CreateCommission(ctx context.Context) (*npool.Commission, error) {
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.getUser(ctx); err != nil {
		return nil, err
	}
	if err := handler.validateGood(ctx); err != nil {
		return nil, err
	}
	if err := handler.createCommmission(ctx); err != nil {
		return nil, err
	}
	info, err := h.GetCommission(ctx)
	if err != nil {
		return nil, err
	}
	handler.notifyCreateCommission(ctx)
	return info, nil
}

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
		Channel:   basetypes.NotifChannel_ChannelEmail,
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
