package invitationcode

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	invitationcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"
	invitationcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"
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
	if len(infos) == 0 {
		return nil, 0, nil
	}

	userIDs := []string{}
	for _, info := range infos {
		userIDs = append(userIDs, info.UserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, 0, int32(len(userIDs)))
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
			Disabled:       info.Disabled,
			CreatedAt:      info.CreatedAt,
			UpdatedAt:      info.UpdatedAt,
		})
	}

	return ivcs, total, nil
}
