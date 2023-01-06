package commission

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"

	comm1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetCommissions(ctx context.Context, in *npool.GetCommissionsRequest) (*npool.GetCommissionsResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.GetCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		return &npool.GetCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := constant.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := comm1.GetCommissions(ctx, &commmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		UserID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetUserID(),
		},
		SettleType: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(in.GetSettleType()),
		},
	}, in.GetOffset(), limit)
	if err != nil {
		return &npool.GetCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCommissionsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppCommissions(ctx context.Context, in *npool.GetAppCommissionsRequest) (*npool.GetAppCommissionsResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.GetAppCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	limit := constant.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := comm1.GetCommissions(ctx, &commmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		SettleType: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(in.GetSettleType()),
		},
	}, in.GetOffset(), limit)
	if err != nil {
		return &npool.GetAppCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCommissionsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
