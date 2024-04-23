package commission

import (
	"context"

	commissionmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	commissionmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
)

type updateHandler struct {
	*Handler
	*checkHandler
}

func (h *Handler) UpdateCommission(ctx context.Context) (*npool.Commission, error) {
	handler := &updateHandler{
		checkHandler: &checkHandler{
			Handler: h,
		},
	}

	if err := handler.checkCommission(ctx); err != nil {
		return nil, err
	}

	_, err := commissionmwcli.UpdateCommission(ctx, &commissionmwpb.CommissionReq{
		ID:        h.ID,
		StartAt:   h.StartAt,
		Threshold: h.Threshold,
	})
	if err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
