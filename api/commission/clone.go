package commission

import (
	"context"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	comm1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"
)

func (s *Server) CloneCommissions(ctx context.Context, in *npool.CloneCommissionsRequest) (*npool.CloneCommissionsResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetNewGoodID()); err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetOldGoodID()); err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	_, total, err := goodmwcli.GetGoods(ctx, &goodmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetNewGoodID(),
		},
	}, 0, 1)
	if err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}
	if total == 0 {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, "new good not exist")
	}

	_, total, err = commmwcli.GetCommissions(ctx, &commmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetOldGoodID(),
		},
		SettleType: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(mgrpb.SettleType_GoodOrderPercent),
		},
	}, 0, 1)
	if err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}
	if total == 0 {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, "old good commission not exist")
	}

	_, total, err = commmwcli.GetCommissions(ctx, &commmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetNewGoodID(),
		},
		SettleType: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(mgrpb.SettleType_GoodOrderPercent),
		},
	}, 0, 1)
	if err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}
	if total > 0 {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, "new good commission already exist")
	}

	err = comm1.CloneCommissions(ctx, in.GetAppID(), in.GetOldGoodID(), in.GetNewGoodID())
	if err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.CloneCommissionsResponse{}, nil
}
