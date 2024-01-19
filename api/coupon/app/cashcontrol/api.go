package cashcontrol

import (
	"context"

	CashControl "github.com/NpoolPlatform/message/npool/inspire/gw/v1/coupon/app/cashcontrol"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	CashControl.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	CashControl.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := CashControl.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
