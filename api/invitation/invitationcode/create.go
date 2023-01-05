package invitationcode

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"
	invitationcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	invitationcode1 "github.com/NpoolPlatform/inspire-gateway/pkg/invitation/invitationcode"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

//nolint
func (s *Server) CreateInvitationCode(ctx context.Context, in *npool.CreateInvitationCodeRequest) (*npool.CreateInvitationCodeResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.CreateInvitationCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		return &npool.CreateInvitationCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := usermwcli.GetUser(ctx, in.GetAppID(), in.GetTargetUserID())
	if err != nil {
		return &npool.CreateInvitationCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if user == nil {
		return &npool.CreateInvitationCodeResponse{}, status.Error(codes.InvalidArgument, "TargetUserID is invalid")
	}

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return &npool.CreateInvitationCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.CreateInvitationCodeResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	info, err := invitationcode1.CreateInvitationCode(ctx, &invitationcodemgrpb.InvitationCodeReq{
		AppID:  &in.AppID,
		UserID: &in.TargetUserID,
	})
	if err != nil {
		return &npool.CreateInvitationCodeResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateInvitationCodeResponse{
		Info: info,
	}, nil
}
