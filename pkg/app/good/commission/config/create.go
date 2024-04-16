package config

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/good/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/good/commission/config"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/good/commission/config"
)

type createHandler struct {
	*Handler
	req    *commissionconfigmwpb.AppGoodCommissionConfigReq
	goodID *string
}

func (h *createHandler) checkGood(ctx context.Context) error {
	if h.AppGoodID == nil {
		return nil
	}

	appgood, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
	})
	if err != nil {
		return err
	}
	if appgood == nil {
		return fmt.Errorf("invalid appgood")
	}

	h.goodID = &appgood.GoodID

	good, err := goodmwcli.GetGoodOnly(ctx, &goodmwpb.Conds{
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: appgood.GoodID},
	})
	if err != nil {
		return err
	}
	if good == nil {
		return fmt.Errorf("invalid good")
	}

	coin, err := appcoinmwcli.GetCoinOnly(ctx, &appcoinmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: appgood.CoinTypeID},
	})
	if err != nil {
		return err
	}
	if coin == nil {
		return fmt.Errorf("invalid coin")
	}
	return nil
}

func (h *createHandler) createCommissionConfig(ctx context.Context) error {
	h.req = &commissionconfigmwpb.AppGoodCommissionConfigReq{
		EntID:           h.EntID,
		AppID:           h.AppID,
		GoodID:          h.goodID,
		AppGoodID:       h.AppGoodID,
		SettleType:      h.SettleType,
		Invites:         h.Invites,
		StartAt:         h.StartAt,
		AmountOrPercent: h.AmountOrPercent,
		ThresholdAmount: h.ThresholdAmount,
		Disabled:        h.Disabled,
		Level:           h.Level,
	}
	if _, err := commissionconfigmwcli.CreateCommissionConfig(ctx, h.req); err != nil {
		return err
	}
	return nil
}

func (h *Handler) CreateCommissionConfig(ctx context.Context) (*npool.AppGoodCommissionConfig, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	handler := &createHandler{
		Handler: h,
	}

	if err := handler.checkGood(ctx); err != nil {
		return nil, err
	}

	if err := handler.createCommissionConfig(ctx); err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
