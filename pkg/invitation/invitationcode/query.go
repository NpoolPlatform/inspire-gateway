package invitationcode

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"

	invitationcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	invitationcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
)

func GetInvitationCode(ctx context.Context, id string) (*npool.InvitationCode, error) {
	info, err := invitationcodemwcli.GetInvitationCode(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := usermwcli.GetUser(ctx, info.AppID, info.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	return &npool.InvitationCode{
		AppID:          info.AppID,
		UserID:         info.UserID,
		EmailAddress:   user.EmailAddress,
		PhoneNO:        user.PhoneNO,
		Username:       user.Username,
		InvitationCode: info.InvitationCode,
		Confirmed:      info.Confirmed,
		Disabled:       info.Disabled,
		CreatedAt:      info.CreatedAt,
		UpdatedAt:      info.UpdatedAt,
	}, nil
}

func GetInvitationCodes(
	ctx context.Context,
	conds *invitationcodemgrpb.Conds,
	offset, limit int32,
) (
	[]*npool.InvitationCode,
	uint32,
	error,
) {
	infos, total, err := invitationcodemwcli.GetInvitationCodes(ctx, conds, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	userIDs := []string{}
	for _, info := range infos {
		userIDs = append(userIDs, info.UserID)
	}

	users, _, err := usermwcli.GetManyUsers(ctx, userIDs)
	if err != nil {
		return nil, 0, err
	}

	userMap := map[string]*usermwpb.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	ivcs := []*npool.InvitationCode{}
	for _, info := range infos {
		user, ok := userMap[info.UserID]
		if !ok {
			continue
		}

		ivcs = append(ivcs, &npool.InvitationCode{
			AppID:          info.AppID,
			UserID:         info.UserID,
			EmailAddress:   user.EmailAddress,
			PhoneNO:        user.PhoneNO,
			Username:       user.Username,
			InvitationCode: info.InvitationCode,
			Confirmed:      info.Confirmed,
			Disabled:       info.Disabled,
			CreatedAt:      info.CreatedAt,
			UpdatedAt:      info.UpdatedAt,
		})
	}

	return ivcs, total, nil
}
