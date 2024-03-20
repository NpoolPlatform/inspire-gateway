package config

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
)

type updateHandler struct {
	*Handler
	info *appconfigmwpb.AppConfig
}

func (h *updateHandler) validateAppConfigs(ctx context.Context) error {
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

func (h *Handler) UpdateAppConfig(ctx context.Context) (*npool.AppConfig, error) {
	info, err := appconfigmwcli.GetAppConfig(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid appconfig")
	}
	if info.ID != *h.ID || info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	handler := &updateHandler{
		Handler: h,
		info:    info,
	}
	if err := handler.validateAppConfigs(ctx); err != nil {
		return nil, err
	}

	_, err = appconfigmwcli.UpdateAppConfig(ctx, &appconfigmwpb.AppConfigReq{
		ID:      h.ID,
		StartAt: h.StartAt,
	})
	if err != nil {
		return nil, err
	}

	return h.GetAppConfig(ctx)
}
