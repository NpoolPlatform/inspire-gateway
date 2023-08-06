package invitationcode

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	invitationcode1 "github.com/NpoolPlatform/inspire-gateway/pkg/invitation/invitationcode"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetInvitationCodes(ctx context.Context, in *npool.GetInvitationCodesRequest) (*npool.GetInvitationCodesResponse, error) {
	handler, err := invitationcode1.NewHandler(
		ctx,
		invitationcode1.WithAppID(&in.AppID),
		invitationcode1.WithOffset(in.GetOffset()),
		invitationcode1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetInvitationCodes",
			"In", in,
			"Err", err,
		)
		return &npool.GetInvitationCodesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetInvitationCodes(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetInvitationCodes",
			"In", in,
			"Err", err,
		)
		return &npool.GetInvitationCodesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetInvitationCodesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppInvitationCodes(ctx context.Context, in *npool.GetAppInvitationCodesRequest) (*npool.GetAppInvitationCodesResponse, error) {
	handler, err := invitationcode1.NewHandler(
		ctx,
		invitationcode1.WithAppID(&in.TargetAppID),
		invitationcode1.WithOffset(in.GetOffset()),
		invitationcode1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppInvitationCodes",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppInvitationCodesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetInvitationCodes(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppInvitationCodes",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppInvitationCodesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppInvitationCodesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
