package archivement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	goodscli "github.com/NpoolPlatform/good-middleware/pkg/client/good"

	"github.com/shopspring/decimal"

	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	userpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	archivementdetailmgrcli "github.com/NpoolPlatform/inspire-manager/pkg/client/archivement/detail"
	archivementgeneralmgrcli "github.com/NpoolPlatform/inspire-manager/pkg/client/archivement/general"
	archivementdetailmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/archivement/detail"
	archivementgeneralmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/archivement/general"

	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	regmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"

	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"

	goodspb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"

	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/good"

	coininfopb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"
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

	invitations, _, err := regmwcli.GetSubordinates(ctx, &regmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		InviterIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: []string{userID},
		},
	}, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	ivMap := map[string]*regmgrpb.Registration{}
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

	invitations, total, err := regmwcli.GetSuperiores(ctx, &regmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		InviteeIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: userIDs,
		},
	}, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if uint32(offset) > total {
		if len(invitations) == 0 {
			return []*npool.UserArchivement{}, 0, nil
		}
	}

	ivMap := map[string]*regmgrpb.Registration{}
	for _, iv := range invitations {
		ivMap[iv.InviteeID] = iv
	}

	return getUserArchivements(ctx, appID, uuid.UUID{}.String(), userIDs, ivMap, offset, limit)
}

// nolint
func getUserArchivements(
	ctx context.Context,
	appID, userID string, uids []string,
	ivMap map[string]*regmgrpb.Registration,
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
		ivs, _, err := regmwcli.GetSubordinates(ctx, &regmgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
			InviterIDs: &commonpb.StringSliceVal{
				Op:    cruder.IN,
				Value: uids,
			},
		}, inviteesOfs, limit)
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
	ofs := 0
	lim := 1000
	coins := []*coininfopb.Coin{}
	for {
		coinInfos, _, err := coininfocli.GetCoins(ctx, nil, int32(ofs), int32(lim))
		if err != nil {
			return nil, 0, err
		}
		if len(coinInfos) == 0 {
			break
		}
		coins = append(coins, coinInfos...)
		ofs += lim
	}

	coinMap := map[string]*coininfopb.Coin{}
	for _, coin := range coins {
		coinMap[coin.ID] = coin
	}

	users, n, err := usercli.GetManyUsers(ctx, uids)
	if err != nil {
		return nil, 0, err
	}

	userMap := map[string]*userpb.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	percents := []*commmwpb.Commission{}
	iofs := int32(0)

	for {
		p, _, err := commmwcli.GetCommissions(ctx, &commmwpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
			UserIDs: &commonpb.StringSliceVal{
				Op:    cruder.IN,
				Value: uids,
			},
		}, iofs, limit)
		if err != nil {
			return nil, 0, err
		}
		if len(p) == 0 {
			break
		}
		percents = append(percents, p...)
		iofs += limit
	}

	goodIDs := []string{}

	for _, val := range percents {
		goodIDs = append(goodIDs, val.GetGoodID())
	}

	for _, val := range generals {
		goodIDs = append(goodIDs, val.GetGoodID())
	}

	goods, _, err := goodscli.GetGoods(ctx, &goodmgrpb.Conds{
		IDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: goodIDs,
		},
	}, 0, int32(len(goodIDs)))
	if err != nil {
		return nil, 0, err
	}

	goodMap := map[string]*goodspb.Good{}
	for _, good := range goods {
		goodMap[good.ID] = good
	}

	archGoodMap := map[string]*goodspb.Good{}

	for _, p := range percents {
		if p.GetGoodID() == "" || p.GetGoodID() == uuid1.InvalidUUIDStr {
			continue
		}
		good, ok := goodMap[p.GetGoodID()]
		if !ok {
			logger.Sugar().Warn("good not exist continue")
			continue
		}

		archGoodMap[p.GetGoodID()] = good
	}

	// 5 Merge info
	archivements := map[string]*npool.UserArchivement{}
	for _, user := range users {
		invitedAt := uint32(0)

		iv, ok := ivMap[user.ID]
		if ok {
			invitedAt = iv.CreatedAt
		}

		archivements[user.ID] = &npool.UserArchivement{
			UserID:        user.ID,
			Username:      user.Username,
			EmailAddress:  user.EmailAddress,
			PhoneNO:       user.PhoneNO,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Kol:           user.Kol,
			TotalInvitees: inviteesMap[user.ID],
			CreatedAt:     user.CreatedAt,
			InvitedAt:     invitedAt,
		}
	}

	for _, general := range generals {
		archivement, ok := archivements[general.UserID]
		if !ok {
			logger.Sugar().Warn("user not exist continue")
			continue
		}

		coin, ok := coinMap[general.CoinTypeID]
		if !ok {
			logger.Sugar().Warn("coin not exist continue")
			continue
		}

		good, ok := goodMap[general.GoodID]
		if !ok {
			logger.Sugar().Warn("good not exist continue")
			continue
		}

		percent := decimal.NewFromInt(0)

		for _, p := range percents {
			if general.UserID != p.UserID || general.GoodID != p.GetGoodID() {
				continue
			}
			percent, err = decimal.NewFromString(p.GetPercent())
			if err != nil {
				continue
			}
			break
		}

		arch := &npool.GoodArchivement{
			GoodID:            general.GoodID,
			GoodName:          good.Title,
			GoodUnit:          good.Unit,
			CommissionPercent: percent.String(),
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
	}

	for _, archivement := range archivements {
	nextCoin:
		for goodID, good := range archGoodMap {
			for _, iarch := range archivement.Archivements {
				if iarch.GoodID == goodID {
					continue nextCoin
				}
			}

			percent := decimal.NewFromInt(0)

			for _, p := range percents {
				if archivement.UserID != p.UserID || goodID != p.GetGoodID() {
					continue
				}
				percent, err = decimal.NewFromString(p.GetPercent())
				if err != nil {
					continue
				}
				break
			}

			coin := coinMap[good.CoinTypeID]

			arch := &npool.GoodArchivement{
				GoodID:            goodID,
				GoodName:          good.Title,
				GoodUnit:          good.Unit,
				CommissionPercent: percent.String(),
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
