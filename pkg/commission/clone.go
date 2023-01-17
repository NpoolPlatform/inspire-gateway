package commission

import (
	"context"
	"fmt"

	commmgrcli "github.com/NpoolPlatform/inspire-manager/pkg/client/commission/goodorderpercent"
	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	commmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"
	goodorderpercentmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission/goodorderpercent"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
)

func CloneCommissions(ctx context.Context, appID, fromGoodID, toGoodID string, settleType commmgrpb.SettleType) error {
	switch settleType {
	case commmgrpb.SettleType_GoodOrderPercent:
		return cloneGoodOrderPercent(ctx, appID, fromGoodID, toGoodID)
	case commmgrpb.SettleType_LimitedOrderPercent:
		fallthrough //nolint
	case commmgrpb.SettleType_AmountThreshold:
		fallthrough //nolint
	case commmgrpb.SettleType_NoCommission:
		return fmt.Errorf("not implemented")
	default:
		return fmt.Errorf("unknown settle type")
	}
}

func cloneGoodOrderPercent(ctx context.Context, appID, fromGoodID, toGoodID string) error {
	offset := int32(0)
	limit := int32(1000) //nolint
	for {
		infos, _, err := commmwcli.GetCommissions(ctx, &commmwpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
			GoodID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: fromGoodID,
			},
			SettleType: &commonpb.Int32Val{
				Op:    cruder.EQ,
				Value: int32(commmgrpb.SettleType_GoodOrderPercent),
			},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(infos) == 0 {
			break
		}
		offset += limit
		req := []*goodorderpercentmgrpb.OrderPercentReq{}
		for _, val := range infos {
			req = append(req, &goodorderpercentmgrpb.OrderPercentReq{
				AppID:   &val.AppID,
				UserID:  &val.UserID,
				GoodID:  &toGoodID,
				Percent: val.Percent,
				StartAt: &val.StartAt,
				EndAt:   &val.EndAt,
			})
		}

		_, err = commmgrcli.CreateOrderPercents(ctx, req)
		if err != nil {
			return err
		}
	}
	return nil
}
