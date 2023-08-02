package commission

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
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
		return err
	}
	if good == nil {
		return fmt.Errorf("invalid good")
	}

	coin, err := appcoinmwcli.GetCoinOnly(ctx, &appcoinmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: good.CoinTypeID},
	})
	if err != nil {
		return err
	}
	if coin == nil {
		return fmt.Errorf("invalid coin")
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
		AmountOrPercent: h.AmountOrPercent,
		Threshold:       h.Threshold,
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
				AmountOrPercent: *h.AmountOrPercent,
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
	if err := handler.createCommission(ctx); err != nil {
		return nil, err
	}
	info, err := h.GetCommission(ctx)
	if err != nil {
		return nil, err
	}
	handler.notifyCreateCommission(ctx)
	return info, nil
}
