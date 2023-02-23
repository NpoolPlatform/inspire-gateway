package event

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
)

func expand(ctx context.Context, info *mgrpb.Event) (*npool.Event, error) {
	app, err := appmwcli.GetApp(ctx, info.AppID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("app is invalid")
	}

	return &npool.Event{
		ID:            info.ID,
		AppID:         info.AppID,
		AppName:       app.Name,
		EventType:     info.EventType,
		Credits:       info.Credits,
		CreditsPerUSD: info.CreditsPerUSD,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.UpdatedAt,
	}, nil
}
