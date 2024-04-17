package config

import (
	"context"

	appconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/config"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/config"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
)

type updateHandler struct {
	*checkHandler
}

func (h *Handler) UpdateAppConfig(ctx context.Context) (*npool.AppConfig, error) {
	handler := &updateHandler{
		checkHandler: &checkHandler{
			Handler: h,
		},
	}
	if err := handler.checkConfig(ctx); err != nil {
		return nil, err
	}

	_, err := appconfigmwcli.UpdateAppConfig(ctx, &appconfigmwpb.AppConfigReq{
		ID:      h.ID,
		StartAt: h.StartAt,
	})
	if err != nil {
		return nil, err
	}

	return h.GetAppConfig(ctx)
}
