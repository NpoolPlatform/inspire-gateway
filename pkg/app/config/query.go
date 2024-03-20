package config

import (
	"context"
	"fmt"

	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
)

type queryHandler struct {
	*Handler
	comms []*appconfigmwpb.AppConfig
	infos []*npool.AppConfig
}

func (h *queryHandler) formalize() {
	for _, comm := range h.comms {
		h.infos = append(h.infos, &npool.AppConfig{
			ID:               comm.ID,
			EntID:            comm.EntID,
			AppID:            comm.AppID,
			SettleMode:       comm.SettleMode,
			SettleAmountType: comm.SettleAmountType,
			SettleInterval:   comm.SettleInterval,
			CommissionType:   comm.CommissionType,
			SettleBenefit:    comm.SettleBenefit,
			StartAt:          comm.StartAt,
			EndAt:            comm.EndAt,
			CreatedAt:        comm.CreatedAt,
			UpdatedAt:        comm.UpdatedAt,
		})
	}
}

func (h *Handler) GetAppConfig(ctx context.Context) (*npool.AppConfig, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}

	info, err := appconfigmwcli.GetAppConfig(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler: h,
		comms:   []*appconfigmwpb.AppConfig{info},
		infos:   []*npool.AppConfig{},
	}

	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}

	return handler.infos[0], nil
}

func (h *Handler) GetAppConfigs(ctx context.Context) ([]*npool.AppConfig, uint32, error) {
	handler := &queryHandler{
		Handler: h,
		infos:   []*npool.AppConfig{},
	}

	conds := &appconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.EndAt != nil {
		conds.EndAt = &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.EndAt}
	}
	infos, total, err := appconfigmwcli.GetAppConfigs(ctx, conds, h.Offset, h.Limit)
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
