package archivement

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	userpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	archivementgeneralmgrcli "github.com/NpoolPlatform/archivement-manager/pkg/client/general"
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
)

// nolint
func GetCoinArchivements(
	ctx context.Context, appID, userID string, offset, limit int32,
) (
	infos []*npool.GetCoinArchivementsResponse_Archivement, total uint32, err error,
) {
	if limit == 0 {
		limit = 1000
	}

	// 1 Get all layered users
	invitations, n, err := inspirecli.GetInvitees(ctx, appID, []string{userID}, offset, limit)
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

	inviteesMap := map[string]uint32{}
	inviteesOfs := int32(0)

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
	users, err := usercli.GetManyUsers(ctx, uids)
	if err != nil {
		return nil, 0, err
	}

	userMap := map[string]*userpb.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	percents := []*inspirepb.Percent{}
	iofs := int32(0)

	for {
		p, _, err := inspirecli.GetActivePercents(ctx, appID, uids, iofs, limit)
		if err != nil {
			return nil, 0, err
		}
		if len(p) == 0 {
			break
		}
		percents = append(percents, p...)
		iofs += limit
	}

	goods, err := goodscli.GetGoods(ctx)
	if err != nil {
		return nil, 0, err
	}

	goodMap := map[string]*goodspb.GoodInfo{}
	for _, good := range goods {
		goodMap[good.ID] = good
	}

	// 5 Merge info
	archivements := map[string]*npool.GetCoinArchivementsResponse_Archivement{}
	for _, user := range users {
		user, ok := userMap[user.ID]
		if !ok {
			return nil, 0, fmt.Errorf("invalid user")
		}

		kol := true
		invitedAt := uint32(0)

		if userID != user.ID {
			iv, ok := ivMap[user.ID]
			if !ok {
				return nil, 0, fmt.Errorf("invalid invitee")
			}
			kol = iv.Kol
			invitedAt = iv.CreatedAt
		}

		archivements[user.ID] = &npool.GetCoinArchivementsResponse_Archivement{
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

	archCoinMap := map[string]*coininfopb.CoinInfo{}
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

		var good *goodspb.GoodInfo
		var percent *inspirepb.Percent

		for _, p := range percents {
			if general.UserID != p.UserID || general.CoinTypeID == p.CoinTypeID {
				continue
			}

			if percent == nil || percent.Percent < p.Percent {
				percent = p
			}
		}

		arch := &npool.CoinArchivement{
			CoinTypeID:      coin.ID,
			CoinName:        coin.Name,
			CoinLogo:        coin.Logo,
			CoinUnit:        coin.Unit,
			TotalUnits:      general.TotalUnits,
			SelfUnits:       general.SelfUnits,
			TotalAmount:     general.TotalAmount,
			SelfAmount:      general.SelfAmount,
			TotalCommission: general.TotalCommission,
			SelfCommission:  general.SelfCommission,
		}
		if percent != nil {
			arch.CurPercent = percent.Percent
			arch.CurGoodID = percent.GoodID
			good, ok = goodMap[percent.GoodID]
			if !ok {
				continue
			}
		}

		if good != nil {
			arch.CurGoodName = good.Title
			arch.CurGoodUnit = good.Unit
		}

		archivement.Archivements = append(archivement.Archivements, arch)
		archivements[general.UserID] = archivement

		archCoinMap[general.CoinTypeID] = coin
		if good != nil || archGoodMap[general.CoinTypeID] == nil {
			archGoodMap[general.CoinTypeID] = good
		}
	}

	for _, archivement := range archivements {
	nextCoin:
		for coinTypeID, coin := range archCoinMap {
			for _, iarch := range archivement.Archivements {
				if iarch.CoinTypeID == coinTypeID {
					continue nextCoin
				}
			}

			arch := &npool.CoinArchivement{
				CoinTypeID:      coin.ID,
				CoinName:        coin.Name,
				CoinLogo:        coin.Logo,
				CoinUnit:        coin.Unit,
				TotalAmount:     decimal.NewFromInt(0).String(),
				SelfAmount:      decimal.NewFromInt(0).String(),
				TotalCommission: decimal.NewFromInt(0).String(),
				SelfCommission:  decimal.NewFromInt(0).String(),
			}

			good := archGoodMap[coinTypeID]
			if good != nil {
				arch.CurGoodID = good.ID
				arch.CurGoodName = good.Title
				arch.CurGoodUnit = good.Unit
			}

			archivement.Archivements = append(archivement.Archivements, arch)
		}
	}

	for _, ar := range archivements {
		infos = append(infos, ar)
	}

	return infos, n, nil
}
