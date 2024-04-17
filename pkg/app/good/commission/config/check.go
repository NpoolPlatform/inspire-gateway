package config

import (
	"context"
	"fmt"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/good/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/good/commission/config"
)

type checkHandler struct {
	*Handler
}

func (h *checkHandler) checkConfig(ctx context.Context) error {
	exist, err := commissionconfigmwcli.ExistCommissionConfigConds(ctx, &commissionconfigmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt: &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid appgoodcommissionconfig")
	}
	return nil
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

	h.GoodID = &appgood.GoodID

	good, err := goodmwcli.GetGoodOnly(ctx, &goodmwpb.Conds{
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: appgood.GoodID},
	})
	if err != nil {
		return err
	}
	if good == nil {
		return fmt.Errorf("invalid good")
	}

	existCoin, err := appcoinmwcli.ExistCoinConds(ctx, &appcoinmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: appgood.CoinTypeID},
	})
	if err != nil {
		return err
	}
	if !existCoin {
		return fmt.Errorf("invalid coin")
	}
	return nil
}
