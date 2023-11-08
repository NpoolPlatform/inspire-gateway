//nolint:dupl
package registration

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	registration1 "github.com/NpoolPlatform/inspire-gateway/pkg/invitation/registration"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRegistrations(ctx context.Context, in *npool.GetRegistrationsRequest) (*npool.GetRegistrationsResponse, error) {
	handler, err := registration1.NewHandler(
		ctx,
		registration1.WithAppID(&in.AppID, true),
		registration1.WithOffset(in.GetOffset()),
		registration1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetRegistration",
			"In", in,
			"Err", err,
		)
		return &npool.GetRegistrationsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetRegistrations(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetRegistration",
			"In", in,
			"Err", err,
		)
		return &npool.GetRegistrationsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetRegistrationsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppRegistrations(ctx context.Context, in *npool.GetAppRegistrationsRequest) (*npool.GetAppRegistrationsResponse, error) {
	handler, err := registration1.NewHandler(
		ctx,
		registration1.WithAppID(&in.TargetAppID, true),
		registration1.WithOffset(in.GetOffset()),
		registration1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppRegistration",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppRegistrationsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetRegistrations(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppRegistration",
			"In", in,
			"Err", err,
		)
		return &npool.GetAppRegistrationsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppRegistrationsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
