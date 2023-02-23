package event

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"

	mgrcli "github.com/NpoolPlatform/inspire-manager/pkg/client/event"
)

func GetEvents(ctx context.Context, conds *mgrpb.Conds, offset, limit int32) ([]*npool.Event, uint32, error) {
	infos, total, err := mgrcli.GetEvents(ctx, conds, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	_infos, err := expandMany(ctx, infos)
	if err != nil {
		return nil, 0, err
	}

	return _infos, total, nil
}
