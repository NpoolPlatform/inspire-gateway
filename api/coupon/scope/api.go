package scope

import (
	"context"

	scope "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/scope"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	scope.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	scope.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := scope.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
