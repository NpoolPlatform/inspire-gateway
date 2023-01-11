package invitationcode

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"
	invitationcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	invitationcode1 "github.com/NpoolPlatform/inspire-gateway/pkg/invitation/invitationcode"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetInvitationCodes(ctx context.Context, in *npool.GetInvitationCodesRequest) (*npool.GetInvitationCodesResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.GetInvitationCodesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := constant.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := invitationcode1.GetInvitationCodes(ctx, &invitationcodemgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
	}, in.GetOffset(), limit)
	if err != nil {
		return &npool.GetInvitationCodesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetInvitationCodesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppInvitationCodes(
	ctx context.Context,
	in *npool.GetAppInvitationCodesRequest,
) (
	*npool.GetAppInvitationCodesResponse,
	error,
) {
	resp, err := s.GetInvitationCodes(ctx, &npool.GetInvitationCodesRequest{
		AppID:  in.GetTargetAppID(),
		Offset: in.GetOffset(),
		Limit:  in.GetLimit(),
	})
	if err != nil {
		return &npool.GetAppInvitationCodesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppInvitationCodesResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}
