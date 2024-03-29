package config

import (
	"context"

	"fmt"

	"github.com/google/uuid"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
)

type createHandler struct {
	*Handler
	req *appconfigmwpb.AppConfigReq
}

func (h *createHandler) createAppConfig(ctx context.Context) error {
	h.req = &appconfigmwpb.AppConfigReq{
		EntID:            h.EntID,
		AppID:            h.AppID,
		CommissionType:   h.CommissionType,
		SettleMode:       h.SettleMode,
		SettleAmountType: h.SettleAmountType,
		SettleInterval:   h.SettleInterval,
		StartAt:          h.StartAt,
		SettleBenefit:    h.SettleBenefit,
	}
	if _, err := appconfigmwcli.CreateAppConfig(ctx, h.req); err != nil {
		return err
	}
	return nil
}

func (h *createHandler) validateAppConfigs(ctx context.Context) error {
	if h.StartAt == nil {
		return nil
	}

	appConfigs := []*appconfigmwpb.AppConfig{}
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		_appConfigs, _, err := appconfigmwcli.GetAppConfigs(ctx, &appconfigmwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			EndAt: &basetypes.Uint32Val{Op: cruder.NEQ, Value: 0},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(_appConfigs) == 0 {
			break
		}
		appConfigs = append(appConfigs, _appConfigs...)
		offset += limit
	}
	for _, appConfig := range appConfigs {
		if appConfig.EndAt > *h.StartAt {
			return fmt.Errorf("invalid startat")
		}
	}
	return nil
}

func (h *Handler) CreateAppConfig(ctx context.Context) (*npool.AppConfig, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	handler := &createHandler{
		Handler: h,
	}

	if err := handler.validateAppConfigs(ctx); err != nil {
		return nil, err
	}

	if err := handler.createAppConfig(ctx); err != nil {
		return nil, err
	}

	return h.GetAppConfig(ctx)
}
