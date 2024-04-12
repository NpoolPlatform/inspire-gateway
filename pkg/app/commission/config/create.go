package config

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/commission/config"
	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/commission/config"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/commission/config"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) validateCommissionCount(ctx context.Context) error {
	appConfig, err := appconfigmwcli.GetAppConfigOnly(ctx, &appconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt: &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
	})
	if err != nil {
		return err
	}
	if appConfig == nil {
		return fmt.Errorf("invalid appconfig")
	}

	offset := int32(0)
	limit := int32(appConfig.MaxLevelCount + 1)
	_commissions, _, err := commissionconfigmwcli.GetCommissionConfigs(ctx, &commissionconfigmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt:      &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
		SettleType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.SettleType)},
		Disabled:   &basetypes.BoolVal{Op: cruder.EQ, Value: false},
	}, offset, limit)
	if err != nil {
		return err
	}

	if len(_commissions) > int(appConfig.MaxLevelCount) {
		return fmt.Errorf("invalid max level")
	}

	return nil
}

func (h *createHandler) validateCommissions(ctx context.Context) error {
	if h.StartAt == nil {
		return nil
	}

	exist, err := commissionconfigmwcli.ExistCommissionConfigConds(ctx, &commissionconfigmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt:      &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
		SettleType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.SettleType)},
		Level:      &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.Level},
	})
	if err != nil {
		return err
	}
	if exist {
		now := uint32(time.Now().Unix())
		if *h.StartAt < now {
			return fmt.Errorf("invalid startat")
		}
	}

	commissions := []*commissionconfigmwpb.AppCommissionConfig{}
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		_commissions, _, err := commissionconfigmwcli.GetCommissionConfigs(ctx, &commissionconfigmwpb.Conds{
			AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			EndAt:      &basetypes.Uint32Val{Op: cruder.NEQ, Value: 0},
			SettleType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.SettleType)},
			Level:      &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.Level},
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
	if err := handler.validateCommissionCount(ctx); err != nil {
		return nil, err
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
		Level:           h.Level,
	}); err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
