package registration

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	registration1 "github.com/NpoolPlatform/inspire-gateway/pkg/invitation/registration"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateRegistration(ctx context.Context, in *npool.UpdateRegistrationRequest) (*npool.UpdateRegistrationResponse, error) {
	handler, err := registration1.NewHandler(
		ctx,
		registration1.WithID(&in.ID),
		registration1.WithAppID(&in.AppID),
		registration1.WithInviterID(&in.InviterID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateRegistration",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateRegistration(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateRegistration",
			"In", in,
			"Err", err,
		)
		return &npool.UpdateRegistrationResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateRegistrationResponse{
		Info: info,
	}, nil
}
