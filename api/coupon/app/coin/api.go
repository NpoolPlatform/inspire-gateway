package coin

import (
	"context"

	coin1 "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/coin"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	coin1.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	coin1.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := coin1.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
