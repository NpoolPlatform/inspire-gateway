package registration

import (
	"context"
	"fmt"

	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"
	regmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"
)

func (h *Handler) UpdateRegistration(ctx context.Context) (*npool.Registration, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	_, err := regmwcli.UpdateRegistration(ctx, &regmwpb.RegistrationReq{
		ID:        h.ID,
		InviterID: h.InviterID,
	})
	if err != nil {
		return nil, err
	}

	return h.GetRegistration(ctx)
}
