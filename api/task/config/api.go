package config

import (
	"context"

	config "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	config.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	config.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := config.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
