package registration

import (
	"context"

	registration "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	registration.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	registration.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := registration.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
