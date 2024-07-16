package reward

import (
	"context"

	reward "github.com/NpoolPlatform/message/npool/inspire/gw/v1/user/coin/reward"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	reward.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	reward.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := reward.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
