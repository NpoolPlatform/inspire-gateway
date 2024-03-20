package config

import (
	"context"
	"fmt"

	commconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/commission/config"
	commconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/commission/config"
)

type queryHandler struct {
	*Handler
	comms []*commconfigmwpb.AppCommissionConfig
	infos []*npool.AppCommissionConfig
}

func (h *queryHandler) formalize() {
	for _, comm := range h.comms {
		h.infos = append(h.infos, &npool.AppCommissionConfig{
			ID:              comm.ID,
			EntID:           comm.EntID,
			AppID:           comm.AppID,
			SettleType:      comm.SettleType,
			AmountOrPercent: comm.AmountOrPercent,
			ThresholdAmount: comm.ThresholdAmount,
			StartAt:         comm.StartAt,
			EndAt:           comm.EndAt,
			CreatedAt:       comm.CreatedAt,
			UpdatedAt:       comm.UpdatedAt,
		})
	}
}

func (h *Handler) GetCommission(ctx context.Context) (*npool.AppCommissionConfig, error) {
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
		Handler: h,
		comms:   []*commconfigmwpb.AppCommissionConfig{info},
		infos:   []*npool.AppCommissionConfig{},
	}

	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}

	return handler.infos[0], nil
}

func (h *Handler) GetCommissions(ctx context.Context) ([]*npool.AppCommissionConfig, uint32, error) {
	handler := &queryHandler{
		Handler: h,
		infos:   []*npool.AppCommissionConfig{},
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

	handler.formalize()
	return handler.infos, total, nil
}
