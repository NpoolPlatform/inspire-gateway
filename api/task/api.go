package task

import (
	"context"

	task "github.com/NpoolPlatform/message/npool/inspire/gw/v1/task"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	task.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	task.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := task.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
