//nolint:dupl
package commission

import (
	"context"
	"fmt"

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

func validateClone(ctx context.Context, appID, oldGoodID, newGoodID string) error {
	if _, err := uuid.Parse(appID); err != nil {
		return err
	}
	if _, err := uuid.Parse(newGoodID); err != nil {
		return err
	}
	if _, err := uuid.Parse(oldGoodID); err != nil {
		return err
	}

	_, total, err := goodmwcli.GetGoods(ctx, &goodmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: newGoodID,
		},
	}, 0, 1)
	if err != nil {
		return err
	}
	if total == 0 {
		return fmt.Errorf("new good not exist")
	}

	_, total, err = commmwcli.GetCommissions(ctx, &commmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: oldGoodID,
		},
		SettleType: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(mgrpb.SettleType_GoodOrderPercent),
		},
	}, 0, 1)
	if err != nil {
		return err
	}
	if total == 0 {
		return fmt.Errorf("old good commission not exist")
	}

	_, total, err = commmwcli.GetCommissions(ctx, &commmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: newGoodID,
		},
		SettleType: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(mgrpb.SettleType_GoodOrderPercent),
		},
	}, 0, 1)
	if err != nil {
		return err
	}
	if total > 0 {
		return fmt.Errorf("new good commission already exist")
	}
	return nil
}

func (s *Server) CloneCommissions(ctx context.Context, in *npool.CloneCommissionsRequest) (*npool.CloneCommissionsResponse, error) {
	err := validateClone(ctx, in.GetAppID(), in.GetOldGoodID(), in.GetNewGoodID())
	if err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = comm1.CloneCommissions(ctx, in.GetAppID(), in.GetOldGoodID(), in.GetNewGoodID())
	if err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.CloneCommissionsResponse{}, nil
}

func (s *Server) CloneAppCommissions(
	ctx context.Context,
	in *npool.CloneAppCommissionsRequest,
) (
	*npool.CloneAppCommissionsResponse,
	error,
) {
	err := validateClone(ctx, in.GetTargetAppID(), in.GetOldGoodID(), in.GetNewGoodID())
	if err != nil {
		return &npool.CloneAppCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = comm1.CloneCommissions(ctx, in.GetTargetAppID(), in.GetOldGoodID(), in.GetNewGoodID())
	if err != nil {
		return &npool.CloneAppCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.CloneAppCommissionsResponse{}, nil
}
