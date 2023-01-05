package registration

import (
	"context"

	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"
	// regmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	// usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
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
