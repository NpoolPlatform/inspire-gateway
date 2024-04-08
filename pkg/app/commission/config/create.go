package config

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/commission/config"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/commission/config"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) validateCommissions(ctx context.Context) error {
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
			SettleType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.SettleType)},
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

func (h *Handler) CreateCommissionConfig(ctx context.Context) (*npool.AppCommissionConfig, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	handler := &createHandler{
		Handler: h,
	}
	if err := handler.validateCommissions(ctx); err != nil {
		return nil, err
	}

	if _, err := commissionconfigmwcli.CreateCommissionConfig(ctx, &commissionconfigmwpb.AppCommissionConfigReq{
		EntID:           h.EntID,
		AppID:           h.AppID,
		SettleType:      h.SettleType,
		Invites:         h.Invites,
		StartAt:         h.StartAt,
		AmountOrPercent: h.AmountOrPercent,
		ThresholdAmount: h.ThresholdAmount,
		Disabled:        h.Disabled,
	}); err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
