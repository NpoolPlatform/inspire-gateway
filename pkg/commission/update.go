package commission

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	commissionmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	commissionmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
)

type updateHandler struct {
	*Handler
	info *commissionmwpb.Commission
}

func (h *updateHandler) validateCommissions(ctx context.Context) error {
	if h.StartAt == nil {
		return nil
	}

	commissions := []*commissionmwpb.Commission{}
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		_commissions, _, err := commissionmwcli.GetCommissions(ctx, &commissionmwpb.Conds{
			AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserID:     &basetypes.StringVal{Op: cruder.EQ, Value: h.info.UserID},
			GoodID:     &basetypes.StringVal{Op: cruder.EQ, Value: h.info.GoodID},
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

func (h *Handler) UpdateCommission(ctx context.Context) (*npool.Commission, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appid")
	}

	info, err := commissionmwcli.GetCommission(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid commission")
	}
	if info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	handler := &updateHandler{
		Handler: h,
		info:    info,
	}
	if err := handler.validateCommissions(ctx); err != nil {
		return nil, err
	}

	_, err = commissionmwcli.UpdateCommission(ctx, &commissionmwpb.CommissionReq{
		ID:        h.ID,
		StartAt:   h.StartAt,
		Threshold: h.Threshold,
	})
	if err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
