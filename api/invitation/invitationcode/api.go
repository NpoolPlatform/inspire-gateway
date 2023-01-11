package invitationcode

import (
	"context"

	invitationcode "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	invitationcode.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	invitationcode.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := invitationcode.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
