package coin

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	couponcoinmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/coin"
	couponcoinmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/app/coin"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
	appcoin *appcoinmwpb.Coin
}

func (h *createHandler) getAppCoin(ctx context.Context) error {
	appcoin, err := appcoinmwcli.GetCoinOnly(ctx, &appcoinmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinTypeID},
	})
	if err != nil {
		return err
	}
	if appcoin == nil {
		return fmt.Errorf("appcoin not exist")
	}
	if !appcoin.StableUSD {
		return fmt.Errorf("not stable usd coin")
	}
	h.appcoin = appcoin
	return nil
}

func (h *createHandler) createCouponCoin(ctx context.Context) error {
	if _, err := couponcoinmwcli.CreateCouponCoin(
		ctx,
		&couponcoinmwpb.CouponCoinReq{
			EntID:      h.EntID,
			AppID:      h.AppID,
			CoinTypeID: h.CoinTypeID,
		},
	); err != nil {
		return err
	}
	return nil
}

func (h *Handler) CreateCouponCoin(ctx context.Context) (*npool.CouponCoin, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	handler := &createHandler{
		Handler: h,
	}
	if err := handler.getAppCoin(ctx); err != nil {
		return nil, err
	}
	if err := handler.createCouponCoin(ctx); err != nil {
		return nil, err
	}

	return h.GetCouponCoin(ctx)
}
