package config

import (
	"context"
	"fmt"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	commconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/good/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/good/commission/config"
	commconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/good/commission/config"

	"github.com/google/uuid"
)

type queryHandler struct {
	*Handler
	appGoods map[string]*appgoodmwpb.Good
	comms    []*commconfigmwpb.AppGoodCommissionConfig
	infos    []*npool.AppGoodCommissionConfig
}

func (h *queryHandler) getAppGoods(ctx context.Context) error {
	goodIDs := []string{}
	for _, comm := range h.comms {
		if _, err := uuid.Parse(comm.AppGoodID); err != nil {
			continue
		}
		goodIDs = append(goodIDs, comm.AppGoodID)
	}
	if len(goodIDs) == 0 {
		return nil
	}

	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: h.comms[0].AppID},
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: goodIDs},
	}, int32(0), int32(len(goodIDs)))
	if err != nil {
		return err
	}

	for _, good := range goods {
		h.appGoods[good.EntID] = good
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, comm := range h.comms {
		appGood, ok := h.appGoods[comm.AppGoodID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.AppGoodCommissionConfig{
			ID:              comm.ID,
			EntID:           comm.EntID,
			AppID:           comm.AppID,
			SettleType:      comm.SettleType,
			GoodID:          comm.GoodID,
			GoodName:        appGood.GoodName,
			AppGoodID:       comm.AppGoodID,
			AppGoodName:     appGood.AppGoodName,
			AmountOrPercent: comm.AmountOrPercent,
			ThresholdAmount: comm.ThresholdAmount,
			StartAt:         comm.StartAt,
			EndAt:           comm.EndAt,
			Invites:         comm.Invites,
			Disabled:        comm.Disabled,
			Level:           comm.Level,
			CreatedAt:       comm.CreatedAt,
			UpdatedAt:       comm.UpdatedAt,
		})
	}
}

func (h *Handler) GetCommission(ctx context.Context) (*npool.AppGoodCommissionConfig, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}

	info, err := commconfigmwcli.GetCommissionConfig(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:  h,
		appGoods: map[string]*appgoodmwpb.Good{},
		comms:    []*commconfigmwpb.AppGoodCommissionConfig{info},
		infos:    []*npool.AppGoodCommissionConfig{},
	}

	if err := handler.getAppGoods(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}

	return handler.infos[0], nil
}

func (h *Handler) GetCommissions(ctx context.Context) ([]*npool.AppGoodCommissionConfig, uint32, error) {
	handler := &queryHandler{
		Handler:  h,
		appGoods: map[string]*appgoodmwpb.Good{},
		infos:    []*npool.AppGoodCommissionConfig{},
	}

	conds := &commconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.EndAt != nil {
		conds.EndAt = &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.EndAt}
	}
	infos, total, err := commconfigmwcli.GetCommissionConfigs(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}
	handler.comms = infos

	if err := handler.getAppGoods(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()
	return handler.infos, total, nil
}
