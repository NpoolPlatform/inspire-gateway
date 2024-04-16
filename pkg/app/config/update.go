package config

import (
	"context"
	"fmt"

	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
)

func (h *Handler) UpdateAppConfig(ctx context.Context) (*npool.AppConfig, error) {
	info, err := appconfigmwcli.GetAppConfig(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid appconfig")
	}
	if info.ID != *h.ID || info.AppID != *h.AppID || info.EndAt != 0 {
		return nil, fmt.Errorf("permission denied")
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
