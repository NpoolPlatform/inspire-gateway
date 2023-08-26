package commission

import (
	"context"
	"fmt"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	commissionmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
)

type cloneHandler struct {
	*Handler
}

func (h *cloneHandler) validateGoods(ctx context.Context) error {
	const limit = 2
	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmgrpb.Conds{
		AppID:   &commonpb.StringVal{Op: cruder.EQ, Value: *h.AppID},
		GoodIDs: &commonpb.StringSliceVal{Op: cruder.IN, Value: []string{*h.FromGoodID, *h.ToGoodID}},
	}, int32(0), int32(limit))
	if err != nil {
		return err
	}
	if len(goods) < limit {
		return fmt.Errorf("invalid goodid")
	}
	return nil
}

func (h *Handler) CloneCommissions(ctx context.Context) error {
	if h.AppID == nil {
		return fmt.Errorf("invalid appid")
	}
	if h.FromGoodID == nil {
		return fmt.Errorf("invalid fromgoodid")
	}
	if h.ToGoodID == nil {
		return fmt.Errorf("invalid togoodid")
	}

	handler := &cloneHandler{
		Handler: h,
	}
	if err := handler.validateGoods(ctx); err != nil {
		return err
	}

	scalePercent := "1"
	if h.ScalePercent != nil {
		scalePercent = *h.ScalePercent
	}
	return commissionmwcli.CloneCommissions(
		ctx,
		*h.AppID,
		*h.FromGoodID,
		*h.ToGoodID,
		scalePercent,
	)
}
