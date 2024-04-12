package config

import (
	"context"
	"time"

	"fmt"

	"github.com/google/uuid"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	appcommissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/commission/config"
	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	appgoodcommissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/good/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"
	appcommissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/commission/config"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
	appgoodcommissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/good/commission/config"
)

type createHandler struct {
	*Handler
	req  *appconfigmwpb.AppConfigReq
	info *appconfigmwpb.AppConfig
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
		MaxLevelCount:    h.MaxLevelCount,
	}
	if _, err := appconfigmwcli.CreateAppConfig(ctx, h.req); err != nil {
		return err
	}
	return nil
}

func (h *createHandler) getAppConfig(ctx context.Context) error {
	appConfig, err := appconfigmwcli.GetAppConfigOnly(ctx, &appconfigmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt: &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
	})
	if err != nil {
		return err
	}
	h.info = appConfig
	return nil
}

func (h *createHandler) validateAppConfigs(ctx context.Context) error {
	if h.StartAt == nil {
		return nil
	}

	if h.info != nil {
		now := uint32(time.Now().Unix())
		if *h.StartAt < now {
			return fmt.Errorf("invalid startat")
		}
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

func (h *createHandler) validateMaxLevelCount(ctx context.Context) error {
	if h.info == nil {
		return nil
	}
	if *h.MaxLevelCount >= h.info.MaxLevelCount {
		return nil
	}

	offset := int32(0)
	limit := int32(h.info.MaxLevelCount*2 + 1)
	_appcommissions, _, err := appcommissionconfigmwcli.GetCommissionConfigs(ctx, &appcommissionconfigmwpb.Conds{
		AppID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EndAt:    &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
		Disabled: &basetypes.BoolVal{Op: cruder.EQ, Value: false},
	}, offset, limit)
	if err != nil {
		return err
	}

	appCommissionCountMap := map[string]uint32{}
	for _, appcommission := range _appcommissions {
		count, ok := appCommissionCountMap[appcommission.SettleType.String()]
		if ok {
			appCommissionCountMap[appcommission.SettleType.String()] = count + 1
			continue
		}
		appCommissionCountMap[appcommission.SettleType.String()] = 1
	}

	for _, count := range appCommissionCountMap {
		if count > *h.MaxLevelCount {
			return fmt.Errorf("invalid maxlevelcount")
		}
	}

	appGoodCommissionCountMap := map[string]uint32{}
	limit = constant.DefaultRowLimit
	for {
		_appgoodcommissions, _, err := appgoodcommissionconfigmwcli.GetCommissionConfigs(ctx, &appgoodcommissionconfigmwpb.Conds{
			AppID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			EndAt:    &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
			Disabled: &basetypes.BoolVal{Op: cruder.EQ, Value: false},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(_appgoodcommissions) == 0 {
			break
		}
		for _, appgoodcommission := range _appgoodcommissions {
			key := fmt.Sprintf("%v_%v", appgoodcommission.AppGoodID, appgoodcommission.SettleType.String())
			count, ok := appGoodCommissionCountMap[key]
			if ok {
				appGoodCommissionCountMap[key] = count + 1
				continue
			}
			appGoodCommissionCountMap[key] = 1
		}
		offset += limit
	}
	for _, count := range appGoodCommissionCountMap {
		if count > *h.MaxLevelCount {
			return fmt.Errorf("invalid maxlevelcount")
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

	if err := handler.getAppConfig(ctx); err != nil {
		return nil, err
	}

	if err := handler.validateAppConfigs(ctx); err != nil {
		return nil, err
	}

	if err := handler.validateMaxLevelCount(ctx); err != nil {
		return nil, err
	}

	if err := handler.createAppConfig(ctx); err != nil {
		return nil, err
	}

	return h.GetAppConfig(ctx)
}
