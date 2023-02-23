package event

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"

	mgrcli "github.com/NpoolPlatform/inspire-manager/pkg/client/event"
)

func UpdateEvent(ctx context.Context, in *mgrpb.EventReq) (*npool.Event, error) {
	info, err := mgrcli.UpdateEvent(ctx, in)
	if err != nil {
		return nil, err
	}

	return expand(ctx, info)
}
