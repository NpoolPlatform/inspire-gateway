package config

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/commission/config"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/commission/config"
)

type updateHandler struct {
	*Handler
	info *commissionconfigmwpb.AppCommissionConfig
}

func (h *updateHandler) validateCommissions(ctx context.Context) error {
	if h.StartAt == nil {
		return nil
	}

	commissions := []*commissionconfigmwpb.AppCommissionConfig{}
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		_commissions, _, err := commissionconfigmwcli.GetCommissionConfigs(ctx, &commissionconfigmwpb.Conds{
			AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			EndAt:      &basetypes.Uint32Val{Op: cruder.NEQ, Value: 0},
			SettleType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(h.info.SettleType)},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(_commissions) == 0 {
			break
		}
		commissions = append(commissions, _commissions...)
		offset += limit
	}
	for _, commission := range commissions {
		if commission.EndAt > *h.StartAt {
			return fmt.Errorf("invalid startat")
		}
	}
	return nil
}

func (h *Handler) UpdateCommission(ctx context.Context) (*npool.AppCommissionConfig, error) {
	info, err := commissionconfigmwcli.GetCommissionConfig(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid commission")
	}
	if info.ID != *h.ID || info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	handler := &updateHandler{
		Handler: h,
		info:    info,
	}
	if err := handler.validateCommissions(ctx); err != nil {
		return nil, err
	}

	_, err = commissionconfigmwcli.UpdateCommissionConfig(ctx, &commissionconfigmwpb.AppCommissionConfigReq{
		ID:              h.ID,
		StartAt:         h.StartAt,
		ThresholdAmount: h.ThresholdAmount,
		Invites:         h.Invites,
	})
	if err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
