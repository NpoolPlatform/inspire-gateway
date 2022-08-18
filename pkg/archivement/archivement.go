package archivement

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	userpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	archivementdetailmgrcli "github.com/NpoolPlatform/archivement-manager/pkg/client/detail"
	archivementgeneralmgrcli "github.com/NpoolPlatform/archivement-manager/pkg/client/general"
	archivementdetailmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/archivement/detail"
	archivementgeneralmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/archivement/general"

	inspirecli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation"
	inspirepb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/inspire/invitation"

	coininfocli "github.com/NpoolPlatform/sphinx-coininfo/pkg/client"

	goodscli "github.com/NpoolPlatform/cloud-hashing-goods/pkg/client"
	goodspb "github.com/NpoolPlatform/message/npool/cloud-hashing-goods"

	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/archivement"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	uuid1 "github.com/NpoolPlatform/go-service-framework/pkg/const/uuid"

	"github.com/google/uuid"
)

func GetGoodArchivements(
	ctx context.Context, appID, userID string, offset, limit int32,
) (
	infos []*npool.UserArchivement, total uint32, err error,
) {
	if limit == 0 {
		limit = 1000
	}

	// 1 Get all layered users
	invitations, _, err := inspirecli.GetInvitees(ctx, appID, []string{userID}, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	ivMap := map[string]*inspirepb.Invitation{}
	for _, iv := range invitations {
		ivMap[iv.InviteeID] = iv
	}

	uids := []string{userID}
	for _, iv := range invitations {
		uids = append(uids, iv.InviteeID)
	}

	return getUserArchivements(ctx, appID, userID, uids, ivMap, offset, limit)
}

func GetUserGoodArchivements(
	ctx context.Context, appID string, userIDs []string, offset, limit int32,
) (
	infos []*npool.UserArchivement, total uint32, err error,
) {
	if limit == 0 {
		limit = 1000
	}

	// 1 Get all layered users
	invitations, total, err := inspirecli.GetInviters(ctx, appID, userIDs, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if uint32(offset) > total {
		if len(invitations) == 0 {
			return []*npool.UserArchivement{}, 0, nil
		}
	}

	ivMap := map[string]*inspirepb.Invitation{}
	for _, iv := range invitations {
		ivMap[iv.InviteeID] = iv
	}

	return getUserArchivements(ctx, appID, uuid.UUID{}.String(), userIDs, ivMap, offset, limit)
}

// nolint
func getUserArchivements(
	ctx context.Context,
	appID, userID string, uids []string,
	ivMap map[string]*inspirepb.Invitation,
	offset, limit int32,
) (
	infos []*npool.UserArchivement, total uint32, err error,
) {
	inviteesMap := map[string]uint32{}
	inviteesOfs := int32(0)

	if limit == 0 {
		limit = 1000
	}

	for {
		ivs, _, err := inspirecli.GetInvitees(ctx, appID, uids, inviteesOfs, limit)
		if err != nil {
			return nil, 0, err
		}
		if len(ivs) == 0 {
			break
		}

		for _, iv := range ivs {
			invitees, ok := inviteesMap[iv.InviterID]
			if !ok {
				inviteesMap[iv.InviterID] = 1
				continue
			}
			inviteesMap[iv.InviterID] = invitees + 1
		}

		inviteesOfs += limit
	}

	// 2 Get all users's archivement
	generals, _, err := archivementgeneralmgrcli.GetGenerals(ctx, &archivementgeneralmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UserIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: uids,
		},
	}, 0, limit)
	if err != nil {
		return nil, 0, err
	}

	// 3 Get coin infos
	coins, err := coininfocli.GetCoinInfos(ctx, cruder.NewFilterConds())
	if err != nil {
		return nil, 0, err
	}

	coinMap := map[string]*coininfopb.CoinInfo{}
	for _, coin := range coins {
		coinMap[coin.ID] = coin
	}

	// 4 Get all users
	users, n, err := usercli.GetManyUsers(ctx, uids)
	if err != nil {
		return nil, 0, err
	}

	userMap := map[string]*userpb.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	goods, err := goodscli.GetGoods(ctx)
	if err != nil {
		return nil, 0, err
	}

	goodMap := map[string]*goodspb.GoodInfo{}
	for _, good := range goods {
		goodMap[good.ID] = good
	}

	percents := []*inspirepb.Percent{}
	iofs := int32(0)

	for {
		p, _, err := inspirecli.GetPercents(ctx, appID, uids, true, iofs, limit)
		if err != nil {
			return nil, 0, err
		}
		if len(p) == 0 {
			break
		}
		percents = append(percents, p...)
		iofs += limit
	}

	for _, p := range percents {
		if p.GoodID == "" || p.GoodID == uuid1.InvalidUUIDStr {
			continue
		}
		good, ok := goodMap[p.GoodID]
		if !ok {
			return nil, 0, fmt.Errorf("invalid good: %v", p)
		}
		if p.CoinTypeID == "" || p.CoinTypeID == uuid1.InvalidUUIDStr {
			p.CoinTypeID = good.CoinInfoID
		}
	}

	// 5 Merge info
	archivements := map[string]*npool.UserArchivement{}
	for _, percent := range percents {
		user, ok := userMap[percent.UserID]
		if !ok {
			return nil, 0, fmt.Errorf("invalid user: %v", user.ID)
		}

		kol := user.ID == userID
		invitedAt := uint32(0)

		iv, ok := ivMap[user.ID]
		if ok {
			kol = iv.Kol
			invitedAt = iv.CreatedAt
		}

		archivements[user.ID] = &npool.UserArchivement{
			UserID:        user.ID,
			Username:      user.Username,
			EmailAddress:  user.EmailAddress,
			PhoneNO:       user.PhoneNO,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Kol:           kol,
			TotalInvitees: inviteesMap[user.ID],
			CreatedAt:     user.CreatedAt,
			InvitedAt:     invitedAt,
		}
	}

	archGoodMap := map[string]*goodspb.GoodInfo{}

	for _, general := range generals {
		archivement, ok := archivements[general.UserID]
		if !ok {
			return nil, 0, fmt.Errorf("invalid general user")
		}

		coin, ok := coinMap[general.CoinTypeID]
		if !ok {
			return nil, 0, fmt.Errorf("invalid coin")
		}

		good, ok := goodMap[general.GoodID]
		if !ok {
			return nil, 0, fmt.Errorf("invalid good: %v", general)
		}

		percent := uint32(0)

		for _, p := range percents {
			if general.UserID != p.UserID || general.GoodID != p.GoodID {
				continue
			}
			percent = p.Percent
			break
		}

		arch := &npool.GoodArchivement{
			GoodID:            general.GoodID,
			GoodName:          good.Title,
			GoodUnit:          good.Unit,
			CommissionPercent: percent,
			CoinTypeID:        coin.ID,
			CoinName:          coin.Name,
			CoinLogo:          coin.Logo,
			CoinUnit:          coin.Unit,
			TotalUnits:        general.TotalUnits,
			SelfUnits:         general.SelfUnits,
			TotalAmount:       general.TotalAmount,
			SelfAmount:        general.SelfAmount,
			TotalCommission:   general.TotalCommission,
			SelfCommission:    general.SelfCommission,
		}

		archivement.Archivements = append(archivement.Archivements, arch)
		archivements[general.UserID] = archivement

		archGoodMap[general.GoodID] = good
	}

	for _, archivement := range archivements {
	nextCoin:
		for goodID, good := range archGoodMap {
			for _, iarch := range archivement.Archivements {
				if iarch.GoodID == goodID {
					continue nextCoin
				}
			}

			percent := uint32(0)

			for _, p := range percents {
				if archivement.UserID != p.UserID || goodID != p.GoodID {
					continue
				}
				percent = p.Percent
				break
			}

			coin := coinMap[good.CoinInfoID]

			arch := &npool.GoodArchivement{
				GoodID:            goodID,
				GoodName:          good.Title,
				GoodUnit:          good.Unit,
				CommissionPercent: percent,
				CoinTypeID:        coin.ID,
				CoinName:          coin.Name,
				CoinLogo:          coin.Logo,
				CoinUnit:          coin.Unit,
				TotalAmount:       decimal.NewFromInt(0).String(),
				SelfAmount:        decimal.NewFromInt(0).String(),
				TotalCommission:   decimal.NewFromInt(0).String(),
				SelfCommission:    decimal.NewFromInt(0).String(),
			}

			archivement.Archivements = append(archivement.Archivements, arch)
		}
	}

	// 6 Get my details for invitees' contribution
	details := []*archivementdetailmgrpb.Detail{}
	detailOfs := int32(0)

	for {
		ds, _, err := archivementdetailmgrcli.GetDetails(ctx, &archivementdetailmgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
			UserID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: userID,
			},
		}, detailOfs, limit)
		if err != nil {
			return nil, 0, err
		}
		if len(ds) == 0 {
			break
		}

		details = append(details, ds...)

		detailOfs += limit
	}

	for _, detail := range details {
		arch, ok := archivements[detail.DirectContributorID]
		if !ok {
			continue
		}

		for _, ar := range arch.Archivements {
			if ar.GoodID != detail.GoodID {
				continue
			}

			amount, _ := decimal.NewFromString(detail.Commission)
			src, _ := decimal.NewFromString(ar.SuperiorCommission)

			ar.SuperiorCommission = amount.Add(src).String()
			break
		}
	}

	for _, ar := range archivements {
		infos = append(infos, ar)
	}

	return infos, n, nil
}
