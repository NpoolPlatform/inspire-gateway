package config

import (
	"context"
	"fmt"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	commissionconfigmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/app/good/commission/config"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	commissionconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/good/commission/config"
)

type cloneHandler struct {
	*Handler
	fromGoodID string
	toGoodID   string
}

func (h *cloneHandler) checkGoods(ctx context.Context) error {
	const limit = 2
	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: []string{
			*h.FromAppGoodID,
			*h.ToAppGoodID,
		}},
	}, int32(0), int32(limit))
	if err != nil {
		return err
	}
	if len(goods) < limit {
		return fmt.Errorf("invalid goodid")
	}
	for _, good := range goods {
		switch good.EntID {
		case *h.FromAppGoodID:
			h.fromGoodID = good.GoodID
		case *h.ToAppGoodID:
			h.toGoodID = good.GoodID
		}
	}
	return nil
}

func (h *Handler) CloneCommissions(ctx context.Context) error {
	handler := &cloneHandler{
		Handler: h,
	}
	if err := handler.checkGoods(ctx); err != nil {
		return err
	}

	scalePercent := "1"
	if h.ScalePercent != nil {
		scalePercent = *h.ScalePercent
	}
	return commissionconfigmwcli.CloneCommissionConfigs(ctx, &commissionconfigmwpb.CloneAppGoodCommissionConfigsRequest{
		AppID:         *h.AppID,
		FromGoodID:    handler.fromGoodID,
		FromAppGoodID: *h.FromAppGoodID,
		ToGoodID:      handler.toGoodID,
		ToAppGoodID:   *h.ToAppGoodID,
		ScalePercent:  scalePercent,
	})
}
