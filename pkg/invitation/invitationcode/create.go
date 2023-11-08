package invitationcode

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"
	"github.com/google/uuid"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	invitationcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	invitationcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"
)

func (h *Handler) CreateInvitationCode(ctx context.Context) (*npool.InvitationCode, error) {
	exist, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid userid")
	}

	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}
	if _, err := invitationcodemwcli.CreateInvitationCode(ctx, &invitationcodemwpb.InvitationCodeReq{
		EntID:  h.EntID,
		AppID:  h.AppID,
		UserID: h.UserID,
	}); err != nil {
		return nil, err
	}
	return h.GetInvitationCode(ctx)
}
