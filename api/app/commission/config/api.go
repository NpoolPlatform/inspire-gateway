package commission

import (
	"context"

	config1 "github.com/NpoolPlatform/message/npool/inspire/gw/v1/app/commission/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	config1.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	config1.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := config1.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
