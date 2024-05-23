package config

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	coinconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/config"
	coinconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) checkAppCoin(ctx context.Context) error {
	exist, err := appcoinmwcli.ExistCoinConds(ctx, &appcoinmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinTypeID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid appcoin")
	}
	return nil
}

func (h *Handler) CreateCoinConfig(ctx context.Context) (*npool.CoinConfig, error) {
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.checkAppCoin(ctx); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	if err := coinconfigmwcli.CreateCoinConfig(ctx, &coinconfigmwpb.CoinConfigReq{
		EntID:      h.EntID,
		AppID:      h.AppID,
		CoinTypeID: h.CoinTypeID,
		MaxValue:   h.MaxValue,
		Allocated:  h.Allocated,
	}); err != nil {
		return nil, err
	}

	return h.GetCoinConfig(ctx)
}
