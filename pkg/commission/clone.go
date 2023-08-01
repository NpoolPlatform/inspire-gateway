package commission

import (
	"context"
	"fmt"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
)

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
	scalePercent := "1"
	if h.ScalePercent != nil {
		scalePercent = *h.ScalePercent
	}
	return commmwcli.CloneCommissions(
		ctx,
		*h.AppID,
		*h.FromGoodID,
		*h.ToGoodID,
		scalePercent,
	)
}
