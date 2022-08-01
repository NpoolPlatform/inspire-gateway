package archivement

import (
	"context"
	"fmt"

	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	userpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	archivementgeneralmgrcli "github.com/NpoolPlatform/archivement-manager/pkg/client/general"
	archivementgeneralmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/archivement/general"

	inspirecli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation"
	inspirepb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/inspire/invitation"

	coininfocli "github.com/NpoolPlatform/sphinx-coininfo/pkg/client"

	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/archivement"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
)

// nolint
func GetCoinArchivements(
	ctx context.Context, appID, userID string, offset, limit int32,
) (
	archivements []*npool.CoinArchivement, total uint32, err error,
) {
	if limit == 0 {
		limit = 1000
	}

	// 1 Get all layered users
	invitations, n, err := inspirecli.GetInvitees(ctx, appID, []string{userID}, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(invitations) == 0 {
		return []*npool.CoinArchivement{}, n, nil
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

	// 5 Merge info
	archivements = []*npool.CoinArchivement{}
	for _, general := range generals {
		user, ok := userMap[general.UserID]
		if !ok {
			return nil, 0, fmt.Errorf("invalid user")
		}

		coin, ok := coinMap[general.CoinTypeID]
		if !ok {
			return nil, 0, fmt.Errorf("invalid coin")
		}

		iv, ok := ivMap[general.UserID]
		if !ok {
			return nil, 0, fmt.Errorf("invalid invitee")
		}

		userPercent := uint32(0)
		for _, percent := range percents {
			if general.UserID == percent.UserID && general.CoinTypeID == percent.CoinTypeID {
				userPercent = percent.Percent
				break
			}
		}

		archivements = append(archivements, &npool.CoinArchivement{
			UserID:       user.ID,
			Username:     user.Username,
			EmailAddress: user.EmailAddress,
			PhoneNO:      user.PhoneNO,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Kol:          iv.Kol,

			TotalInvitees: inviteesMap[user.ID],

			CoinTypeID: coin.ID,
			CoinName:   coin.Name,
			CoinLogo:   coin.Logo,
			CoinUnit:   coin.Unit,

			TotalUnits:      general.TotalUnits,
			SelfUnits:       general.SelfUnits,
			TotalAmount:     general.TotalAmount,
			SelfAmount:      general.SelfAmount,
			TotalCommission: general.TotalCommission,
			SelfCommission:  general.SelfCommission,

			CurPercent: userPercent,
		})
	}

	return archivements, n, nil
}
