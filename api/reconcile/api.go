package reconcile

import (
	"context"

	reconcile "github.com/NpoolPlatform/message/npool/inspire/gw/v1/reconcile"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	reconcile.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	reconcile.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := reconcile.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
