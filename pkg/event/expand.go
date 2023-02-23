package event

import (
	"context"
	"fmt"

	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
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

func expandMany(ctx context.Context, infos []*mgrpb.Event) ([]*npool.Event, error) {
	ids := []string{}
	for _, info := range infos {
		ids = append(ids, info.AppID)
	}

	apps, _, err := appmwcli.GetManyApps(ctx, ids)
	if err != nil {
		return nil, err
	}

	appMap := map[string]*appmwpb.App{}
	for _, app := range apps {
		appMap[app.ID] = app
	}

	_infos := []*npool.Event{}

	for _, info := range infos {
		app, ok := appMap[info.AppID]
		if !ok {
			continue
		}

		_infos = append(_infos, &npool.Event{
			ID:            info.ID,
			AppID:         info.AppID,
			AppName:       app.Name,
			EventType:     info.EventType,
			Credits:       info.Credits,
			CreditsPerUSD: info.CreditsPerUSD,
			CreatedAt:     info.CreatedAt,
			UpdatedAt:     info.UpdatedAt,
		})
	}

	return _infos, nil
}
