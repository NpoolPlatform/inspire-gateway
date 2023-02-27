package event

import (
	"context"
	"fmt"

	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/appgood"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	allocatedmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"
	coupmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/coupon"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	coupmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/coupon"

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
		InviterLayers:  info.InviterLayers,
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

	coupons := map[allocatedmgrpb.CouponType][]*coupmwpb.Coupon{}
	couponIDs := map[allocatedmgrpb.CouponType][]string{}
	for _, coup := range info.Coupons {
		couponIDs[coup.CouponType] = append(couponIDs[coup.CouponType], coup.ID)
	}

	for ct, ids := range couponIDs {
		_coupons, _, err := coupmwcli.GetCoupons(ctx, &coupmwpb.Conds{
			CouponType: &commonpb.Int32Val{Op: cruder.EQ, Value: int32(ct)},
			IDs:        &commonpb.StringSliceVal{Op: cruder.IN, Value: ids},
		}, 0, int32(len(ids)))
		if err != nil {
			return nil, err
		}

		coupons[ct] = _coupons
	}

	coupMap := map[string]*coupmwpb.Coupon{}
	for _, coups := range coupons {
		for _, coup := range coups {
			coupMap[coup.ID] = coup
		}
	}

	for _, coup := range info.Coupons {
		_coup, ok := coupMap[coup.ID]
		if !ok {
			continue
		}

		_info.Coupons = append(_info.Coupons, &npool.Coupon{
			ID:         coup.ID,
			CouponType: coup.CouponType,
			Value:      _coup.Value,
			Name:       _coup.Name,
		})
	}

	return _info, nil
}

func expandMany(ctx context.Context, infos []*mgrpb.Event) ([]*npool.Event, error) { //nolint
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

	coupons := map[allocatedmgrpb.CouponType][]*coupmwpb.Coupon{}
	couponIDs := map[allocatedmgrpb.CouponType][]string{}
	for _, info := range infos {
		for _, coup := range info.Coupons {
			couponIDs[coup.CouponType] = append(couponIDs[coup.CouponType], coup.ID)
		}
	}

	for ct, ids := range couponIDs {
		_coupons, _, err := coupmwcli.GetCoupons(ctx, &coupmwpb.Conds{
			CouponType: &commonpb.Int32Val{Op: cruder.EQ, Value: int32(ct)},
			IDs:        &commonpb.StringSliceVal{Op: cruder.IN, Value: ids},
		}, 0, int32(len(ids)))
		if err != nil {
			return nil, err
		}

		coupons[ct] = _coupons
	}

	coupMap := map[string]*coupmwpb.Coupon{}
	for _, coups := range coupons {
		for _, coup := range coups {
			coupMap[coup.ID] = coup
		}
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
			InviterLayers:  info.InviterLayers,
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

		for _, coup := range info.Coupons {
			_coup, ok := coupMap[coup.ID]
			if !ok {
				continue
			}

			_info.Coupons = append(_info.Coupons, &npool.Coupon{
				ID:         coup.ID,
				CouponType: coup.CouponType,
				Value:      _coup.Value,
				Name:       _coup.Name,
			})
		}

		_infos = append(_infos, _info)
	}

	return _infos, nil
}
