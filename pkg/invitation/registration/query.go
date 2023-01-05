package registration

import (
	"context"

	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"
	regmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
)

func GetRegistration(ctx context.Context, id string) (*npool.Registration, error) {
	info, err := regmwcli.GetRegistration(ctx, id)
	if err != nil {
		return nil, err
	}

	userIDs := []string{info.InviterID, info.InviteeID}

	users, _, err := usermwcli.GetManyUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	reg := &npool.Registration{
		AppID:     info.AppID,
		InviterID: info.InviterID,
		InviteeID: info.InviteeID,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}

	for _, user := range users {
		if info.InviterID == user.ID {
			reg.InviterEmailAddress = user.EmailAddress
			reg.InviterPhoneNO = user.PhoneNO
			reg.InviterUsername = user.Username
		}
		if info.InviteeID == user.ID {
			reg.InviteeEmailAddress = user.EmailAddress
			reg.InviteePhoneNO = user.PhoneNO
			reg.InviteeUsername = user.Username
		}
	}

	return reg, nil
}

func GetRegistrations(ctx context.Context, conds *regmgrpb.Conds, offset, limit int32) ([]*npool.Registration, uint32, error) {
	infos, total, err := regmwcli.GetRegistrations(ctx, conds, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	userIDs := []string{}
	for _, info := range infos {
		userIDs = append(userIDs, info.InviterID, info.InviteeID)
	}

	users, _, err := usermwcli.GetManyUsers(ctx, userIDs)
	if err != nil {
		return nil, 0, err
	}

	userMap := map[string]*usermwpb.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	regs := []*npool.Registration{}

	for _, info := range infos {
		inviter, ok := userMap[info.InviterID]
		if !ok {
			continue
		}

		invitee, ok := userMap[info.InviteeID]
		if !ok {
			continue
		}

		regs = append(regs, &npool.Registration{
			AppID:               info.AppID,
			InviterID:           info.InviterID,
			InviterEmailAddress: inviter.EmailAddress,
			InviterPhoneNO:      inviter.PhoneNO,
			InviterUsername:     inviter.Username,
			InviteeID:           info.InviteeID,
			InviteeEmailAddress: invitee.EmailAddress,
			InviteePhoneNO:      invitee.PhoneNO,
			InviteeUsername:     invitee.Username,
			CreatedAt:           info.CreatedAt,
			UpdatedAt:           info.UpdatedAt,
		})
	}

	return regs, total, nil
}
