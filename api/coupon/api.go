package coupon

import (
	"context"

	coupon "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	coupon.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	coupon.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := coupon.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
