package event

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"
)

func CreateEvent(ctx context.Context, in *mgrpb.EventReq) (*npool.Event, error) {
	return nil, nil
}
