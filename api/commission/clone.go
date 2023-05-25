//nolint:dupl
package commission

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"

	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	comm1 "github.com/NpoolPlatform/inspire-gateway/pkg/commission"
)

func validateClone(ctx context.Context, appID, fromGoodID, toGoodID string) error {
	if _, err := uuid.Parse(appID); err != nil {
		return err
	}
	if _, err := uuid.Parse(toGoodID); err != nil {
		return err
	}
	if _, err := uuid.Parse(fromGoodID); err != nil {
		return err
	}

	ag, err := goodmwcli.GetGoodOnly(ctx, &goodmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: fromGoodID,
		},
	})
	if err != nil {
		return err
	}
	if ag == nil {
		return fmt.Errorf("invalid appgood")
	}

	ag, err = goodmwcli.GetGoodOnly(ctx, &goodmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: toGoodID,
		},
	})
	if err != nil {
		return err
	}
	if ag == nil {
		return fmt.Errorf("invalid appgood")
	}

	return nil
}

func (s *Server) CloneCommissions(ctx context.Context, in *npool.CloneCommissionsRequest) (*npool.CloneCommissionsResponse, error) {
	err := validateClone(ctx, in.GetAppID(), in.GetFromGoodID(), in.GetToGoodID())
	if err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = comm1.CloneCommissions(ctx, in.GetAppID(), in.GetFromGoodID(), in.GetToGoodID(), in.GetValue(), in.GetSettleType())
	if err != nil {
		return &npool.CloneCommissionsResponse{}, status.Error(codes.Internal, err.Error())
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
	err := validateClone(ctx, in.GetTargetAppID(), in.GetFromGoodID(), in.GetToGoodID())
	if err != nil {
		return &npool.CloneAppCommissionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = comm1.CloneCommissions(ctx, in.GetTargetAppID(), in.GetFromGoodID(), in.GetToGoodID(), in.GetValue(), in.GetSettleType())
	if err != nil {
		return &npool.CloneAppCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CloneAppCommissionsResponse{}, nil
}
