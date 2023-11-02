package scope

import (
	"context"

	appgoodscope "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/scope"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	appgoodscope.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	appgoodscope.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := appgoodscope.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
