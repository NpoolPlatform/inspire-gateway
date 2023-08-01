package invitationcode

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"

	invitationcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	invitationcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"
)

func (h *Handler) CreateInvitationCode(ctx context.Context) (*npool.InvitationCode, error) {
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
