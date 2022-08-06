package reconciliation

import (
	"context"

	reconciliation "github.com/NpoolPlatform/message/npool/inspire/gw/v1/reconciliation"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	reconciliation.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	reconciliation.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := reconciliation.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
