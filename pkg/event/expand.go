package event

import (
	"context"
	"fmt"

	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/appgood"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	uuid1 "github.com/NpoolPlatform/go-service-framework/pkg/const/uuid"
)

func expand(ctx context.Context, info *mgrpb.Event) (*npool.Event, error) {
	app, err := appmwcli.GetApp(ctx, info.AppID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("app is invalid")
	}

	_info := &npool.Event{
		ID:             info.ID,
		AppID:          info.AppID,
		AppName:        app.Name,
		EventType:      info.EventType,
		Credits:        info.Credits,
		CreditsPerUSD:  info.CreditsPerUSD,
		MaxConsecutive: info.MaxConsecutive,
		CreatedAt:      info.CreatedAt,
		UpdatedAt:      info.UpdatedAt,
	}

	if info.GoodID != uuid1.InvalidUUIDStr {
		good, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmgrpb.Conds{
			AppID:  &commonpb.StringVal{Op: cruder.EQ, Value: info.AppID},
			GoodID: &commonpb.StringVal{Op: cruder.EQ, Value: info.GoodID},
		})
		if err != nil {
			return nil, err
		}
		if good == nil {
			return nil, fmt.Errorf("good is invalid")
		}

		_info.GoodID = info.GoodID
		_info.GoodName = good.GoodName
	}

	return _info, nil
}

func expandMany(ctx context.Context, infos []*mgrpb.Event) ([]*npool.Event, error) {
	appIDs := []string{}
	goodIDs := []string{}

	for _, info := range infos {
		appIDs = append(appIDs, info.AppID)
		if info.GoodID != uuid1.InvalidUUIDStr {
			goodIDs = append(goodIDs, info.GoodID)
		}
	}

	apps, _, err := appmwcli.GetManyApps(ctx, appIDs)
	if err != nil {
		return nil, err
	}

	appMap := map[string]*appmwpb.App{}
	for _, app := range apps {
		appMap[app.ID] = app
	}

	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmgrpb.Conds{
		AppIDs:  &commonpb.StringSliceVal{Op: cruder.IN, Value: appIDs},
		GoodIDs: &commonpb.StringSliceVal{Op: cruder.IN, Value: goodIDs},
	}, 0, int32(len(appIDs)*len(goodIDs)))
	if err != nil {
		return nil, err
	}

	goodMap := map[string]*appgoodmwpb.Good{}
	for _, good := range goods {
		goodMap[good.GoodID] = good
	}

	_infos := []*npool.Event{}

	for _, info := range infos {
		app, ok := appMap[info.AppID]
		if !ok {
			continue
		}

		_info := &npool.Event{
			ID:             info.ID,
			AppID:          info.AppID,
			AppName:        app.Name,
			EventType:      info.EventType,
			Credits:        info.Credits,
			CreditsPerUSD:  info.CreditsPerUSD,
			MaxConsecutive: info.MaxConsecutive,
			GoodID:         info.GoodID,
			CreatedAt:      info.CreatedAt,
			UpdatedAt:      info.UpdatedAt,
		}

		if info.GoodID != uuid1.InvalidUUIDStr {
			good, ok := goodMap[info.GoodID]
			if !ok {
				continue
			}

			_info.GoodID = info.GoodID
			_info.GoodName = good.GoodName
		}

		_infos = append(_infos, _info)
	}

	return _infos, nil
}
