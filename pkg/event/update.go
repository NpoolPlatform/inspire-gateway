package event

import (
	"context"
	"fmt"

	coinconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	coinconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
)

type updateHandler struct {
	*Handler
}

func (h *updateHandler) checkRepeatedCoinConfigs() error {
	if len(h.Coins) == 0 {
		return nil
	}
	coinConfigIDs := map[string]string{}
	for _, coin := range h.Coins {
		_, ok := coinConfigIDs[*coin.CoinConfigID]
		if ok {
			return fmt.Errorf("repeated coinconfig")
		}
		coinConfigIDs[*coin.CoinConfigID] = *coin.CoinConfigID
	}
	return nil
}

//nolint:dupl
func (h *updateHandler) checkCoinConfigs(ctx context.Context) error {
	if len(h.Coins) == 0 {
		return nil
	}
	for _, coin := range h.Coins {
		exist, err := coinconfigmwcli.ExistCoinConfigConds(ctx, &coinconfigmwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *coin.CoinConfigID},
		})
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("invalid coin")
		}
	}

	return nil
}

func (h *Handler) UpdateEvent(ctx context.Context) (*npool.Event, error) {
	info, err := h.GetEvent(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid event")
	}
	if info.ID != *h.ID || info.EntID != *h.EntID || info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	handler := &updateHandler{
		Handler: h,
	}
	if err := handler.checkRepeatedCoinConfigs(); err != nil {
		return nil, err
	}
	if err := handler.checkCoinConfigs(ctx); err != nil {
		return nil, err
	}

	if _, err := eventmwcli.UpdateEvent(ctx, &eventmwpb.EventReq{
		ID:              h.ID,
		EntID:           h.EntID,
		AppID:           h.AppID,
		CouponIDs:       h.CouponIDs,
		Coins:           h.Coins,
		Credits:         h.Credits,
		CreditsPerUSD:   h.CreditsPerUSD,
		MaxConsecutive:  h.MaxConsecutive,
		InviterLayers:   h.InviterLayers,
		RemoveCouponIDs: h.RemoveCouponIDs,
		RemoveCoins:     h.RemoveCoins,
	}); err != nil {
		return nil, err
	}
	return h.GetEvent(ctx)
}
