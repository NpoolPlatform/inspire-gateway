package config

import (
	"context"

	"github.com/google/uuid"

	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/good/commission/config"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/good/commission/config"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/good/commission/config"
)

type createHandler struct {
	*checkHandler
	req *commissionconfigmwpb.AppGoodCommissionConfigReq
}

func (h *createHandler) createCommissionConfig(ctx context.Context) error {
	h.req = &commissionconfigmwpb.AppGoodCommissionConfigReq{
		EntID:           h.EntID,
		AppID:           h.AppID,
		GoodID:          &h.appgood.GoodID,
		AppGoodID:       h.AppGoodID,
		SettleType:      h.SettleType,
		Invites:         h.Invites,
		StartAt:         h.StartAt,
		AmountOrPercent: h.AmountOrPercent,
		ThresholdAmount: h.ThresholdAmount,
		Disabled:        h.Disabled,
		Level:           h.Level,
	}
	if _, err := commissionconfigmwcli.CreateCommissionConfig(ctx, h.req); err != nil {
		return err
	}
	return nil
}

func (h *Handler) CreateCommissionConfig(ctx context.Context) (*npool.AppGoodCommissionConfig, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	handler := &createHandler{
		checkHandler: &checkHandler{
			Handler: h,
		},
	}

	if err := handler.checkGood(ctx); err != nil {
		return nil, err
	}

	if err := handler.createCommissionConfig(ctx); err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
