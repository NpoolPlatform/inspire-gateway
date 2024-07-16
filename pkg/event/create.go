package event

import (
	"context"
	"fmt"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	coinconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coin/config"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	coinconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coin/config"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
	appGood *appgoodmwpb.Good
}

func (h *createHandler) checkRepeatedCoinConfigs() error {
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
func (h *createHandler) checkCoinConfigs(ctx context.Context) error {
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

func (h *createHandler) checkAppGood(ctx context.Context) error {
	if h.AppGoodID == nil {
		return nil
	}

	good, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
	})
	if err != nil {
		return err
	}
	if good == nil {
		return fmt.Errorf("invalid goodid")
	}

	h.appGood = good
	return nil
}

func (h *Handler) CreateEvent(ctx context.Context) (*npool.Event, error) {
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.checkAppGood(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkRepeatedCoinConfigs(); err != nil {
		return nil, err
	}
	if err := handler.checkCoinConfigs(ctx); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	req := &eventmwpb.EventReq{
		EntID:          h.EntID,
		AppID:          h.AppID,
		EventType:      h.EventType,
		CouponIDs:      h.CouponIDs,
		Credits:        h.Credits,
		CreditsPerUSD:  h.CreditsPerUSD,
		MaxConsecutive: h.MaxConsecutive,
		InviterLayers:  h.InviterLayers,
		Coins:          h.Coins,
	}
	if handler.appGood != nil {
		req.GoodID = &handler.appGood.GoodID
		req.AppGoodID = h.AppGoodID
	}

	if _, err := eventmwcli.CreateEvent(ctx, req); err != nil {
		return nil, err
	}

	return h.GetEvent(ctx)
}
