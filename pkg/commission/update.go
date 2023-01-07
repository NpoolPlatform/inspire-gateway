package commission

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	commmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
)

func UpdateCommission(
	ctx context.Context,
	id, appID string,
	settleType mgrpb.SettleType,
	value *string,
	startAt *uint32,
) (
	*npool.Commission,
	error,
) {
	info, err := commmwcli.GetCommission(ctx, id, settleType)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid commission")
	}
	if info.AppID != appID {
		return nil, fmt.Errorf("permission denied")
	}

	req := &commmwpb.CommissionReq{
		ID:         &id,
		SettleType: &settleType,
		StartAt:    startAt,
	}

	if value != nil {
		switch info.SettleType {
		case commmgrpb.SettleType_GoodOrderPercent:
			req.Percent = value
		case commmgrpb.SettleType_LimitedOrderPercent:
			fallthrough //nolint
		case commmgrpb.SettleType_AmountThreshold:
			fallthrough //nolint
		case commmgrpb.SettleType_NoCommission:
			return nil, fmt.Errorf("not implemented")
		default:
			return nil, fmt.Errorf("unknown settle type")
		}
	}

	_, err = commmwcli.UpdateCommission(ctx, req)
	if err != nil {
		return nil, err
	}

	return GetCommission(ctx, id, info.SettleType)
}
