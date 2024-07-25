package coin

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	coin1 "github.com/NpoolPlatform/inspire-gateway/pkg/event/coin"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event/coin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AdminGetEventCoins(ctx context.Context, in *npool.AdminGetEventCoinsRequest) (*npool.AdminGetEventCoinsResponse, error) {
	handler, err := coin1.NewHandler(
		ctx,
		coin1.WithAppID(&in.TargetAppID, true),
		coin1.WithOffset(in.GetOffset()),
		coin1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetEventCoins",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetEventCoinsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetEventCoins(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AdminGetEventCoins",
			"In", in,
			"Err", err,
		)
		return &npool.AdminGetEventCoinsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AdminGetEventCoinsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
