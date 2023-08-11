package achievement

import (
	"context"

	achievement "github.com/NpoolPlatform/message/npool/inspire/gw/v1/achievement"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	achievement.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	achievement.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := achievement.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
