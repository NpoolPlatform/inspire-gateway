package invitationcode

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	invitationcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	invitationcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"
)

func (h *Handler) CreateInvitationCode(ctx context.Context) (*npool.InvitationCode, error) {
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appid")
	}
	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userid")
	}

	exist, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid userid")
	}

	info, err := invitationcodemwcli.CreateInvitationCode(ctx, &invitationcodemwpb.InvitationCodeReq{
		AppID:  h.AppID,
		UserID: h.UserID,
	})
	if err != nil {
		return nil, err
	}
	h.ID = &info.ID
	return h.GetInvitationCode(ctx)
}
