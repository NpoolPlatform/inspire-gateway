package registration

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"
	registrationmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"

	registration1 "github.com/NpoolPlatform/inspire-gateway/pkg/invitation/registration"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetRegistrations(ctx context.Context, in *npool.GetRegistrationsRequest) (*npool.GetRegistrationsResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.GetRegistrationsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := constant.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := registration1.GetRegistrations(ctx, &registrationmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.AppID,
		},
	}, in.GetOffset(), limit)
	if err != nil {
		return &npool.GetRegistrationsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetRegistrationsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
