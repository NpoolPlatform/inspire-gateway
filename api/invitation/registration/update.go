package registration

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"
	registrationmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"

	registration1 "github.com/NpoolPlatform/inspire-gateway/pkg/invitation/registration"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	ivcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) UpdateRegistration(ctx context.Context, in *npool.UpdateRegistrationRequest) (*npool.UpdateRegistrationResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetInviterID()); err != nil {
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := usermwcli.GetUser(ctx, in.GetAppID(), in.GetInviterID())
	if err != nil {
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if user == nil {
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.InvalidArgument, "InviterID is invalid")
	}

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	ivcode, err := ivcodemwcli.GetInvitationCodeOnly(ctx, &ivcodemgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		UserID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetInviterID(),
		},
	})
	if err != nil {
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if ivcode == nil {
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.InvalidArgument, "InviterID is not inviter")
	}

	info, err := registration1.UpdateRegistration(ctx, &registrationmgrpb.RegistrationReq{
		AppID:     &in.AppID,
		InviterID: &in.InviterID,
	})
	if err != nil {
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateRegistrationResponse{
		Info: info,
	}, nil
}
