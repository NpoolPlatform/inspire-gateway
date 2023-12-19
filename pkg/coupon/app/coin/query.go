package coin

import (
	"context"
	"fmt"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	couponcoinmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/app/coin"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/coin"
	couponcoinmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/app/coin"
)

type queryHandler struct {
	*Handler
	infos       []*npool.CouponCoin
	couponcoins []*couponcoinmwpb.CouponCoin
	appcoins    map[string]*appcoinmwpb.Coin
}

func (h *queryHandler) getAppCoins(ctx context.Context) error {
	ids := []string{}
	for _, info := range h.couponcoins {
		ids = append(ids, info.CoinTypeID)
	}

	appcoins, _, err := appcoinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		CoinTypeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, int32(0), int32(len(ids)))
	if err != nil {
		return err
	}
	for _, appcoin := range appcoins {
		h.appcoins[appcoin.EntID] = appcoin
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, info := range h.couponcoins {
		appcoin, ok := h.appcoins[info.CoinTypeID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.CouponCoin{
			ID:                 info.ID,
			EntID:              info.EntID,
			AppID:              info.AppID,
			CouponID:           info.CouponID,
			CouponName:         info.CouponName,
			CouponDenomination: info.CouponDenomination,
			CoinTypeID:         info.CoinTypeID,
			CoinName:           appcoin.CoinName,
			CoinENV:            appcoin.ENV,
			CreatedAt:          info.CreatedAt,
			UpdatedAt:          info.UpdatedAt,
		})
	}
}

func (h *Handler) GetCouponCoin(ctx context.Context) (*npool.CouponCoin, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}
	info, err := couponcoinmwcli.GetCouponCoin(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:     h,
		couponcoins: []*couponcoinmwpb.CouponCoin{info},
		appcoins:    map[string]*appcoinmwpb.Coin{},
	}
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, err
	}

	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetCouponCoins(ctx context.Context) ([]*npool.CouponCoin, uint32, error) {
	couponcoins, total, err := couponcoinmwcli.GetCouponCoins(ctx, &couponcoinmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}

	handler := &queryHandler{
		Handler:     h,
		couponcoins: couponcoins,
		appcoins:    map[string]*appcoinmwpb.Coin{},
	}
	if err := handler.getAppCoins(ctx); err != nil {
		return nil, 0, err
	}

	handler.formalize()
	return handler.infos, total, nil
}