package invitationcode

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"

	invitationcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	invitationcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"
)

func CreateInvitationCode(ctx context.Context, in *invitationcodemgrpb.InvitationCodeReq) (*npool.InvitationCode, error) {
	info, err := invitationcodemwcli.CreateInvitationCode(ctx, in)
	if err != nil {
		return nil, err
	}

	return GetInvitationCode(ctx, info.ID)
}
