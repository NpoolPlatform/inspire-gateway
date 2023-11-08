package registration

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"
	regmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"
)

func (h *Handler) UpdateRegistration(ctx context.Context) (*npool.Registration, error) {
	info, err := h.GetRegistration(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid id")
	}
	if info.ID != *h.ID || info.EntID != *h.EntID {
		return nil, fmt.Errorf("permission denied")
	}
	if info.InviterID == *h.InviterID || info.InviteeID == *h.InviterID {
		return nil, fmt.Errorf("invalid inviterid")
	}

	exist, err := usermwcli.ExistUser(ctx, info.AppID, *h.InviterID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid inviterid")
	}

	_, err = regmwcli.UpdateRegistration(ctx, &regmwpb.RegistrationReq{
		ID:        h.ID,
		AppID:     &info.AppID,
		InviterID: h.InviterID,
	})
	if err != nil {
		return nil, err
	}

	return h.GetRegistration(ctx)
}
