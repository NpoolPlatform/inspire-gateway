package registration

import (
	"context"

	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"
	regmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"
)

func UpdateRegistration(ctx context.Context, in *regmgrpb.RegistrationReq) (*npool.Registration, error) {
	info, err := regmwcli.UpdateRegistration(ctx, in)
	if err != nil {
		return nil, err
	}

	return GetRegistration(ctx, info.ID)
}
