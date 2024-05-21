package allocated

import (
	"context"

	allocated "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coin/allocated"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	allocated.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	allocated.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := allocated.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
