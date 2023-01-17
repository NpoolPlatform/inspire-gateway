package commission

import (
	"context"

	commmgrcli "github.com/NpoolPlatform/inspire-manager/pkg/client/commission/goodorderpercent"
	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"
	commmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission/goodorderpercent"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
)

func CloneCommissions(ctx context.Context, appID, oldGoodID, newGoodID string) error {
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
				Value: oldGoodID,
			},
			SettleType: &commonpb.Int32Val{
				Op:    cruder.EQ,
				Value: int32(mgrpb.SettleType_GoodOrderPercent),
			},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(infos) == 0 {
			break
		}
		offset += limit
		req := []*commmgrpb.OrderPercentReq{}
		for _, val := range infos {
			req = append(req, &commmgrpb.OrderPercentReq{
				AppID:   &val.AppID,
				UserID:  &val.UserID,
				GoodID:  &newGoodID,
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
