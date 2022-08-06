package api

import (
	"context"

	inspire "github.com/NpoolPlatform/message/npool/inspire/gw/v1"

	"github.com/NpoolPlatform/inspire-gateway/api/archivement"
	"github.com/NpoolPlatform/inspire-gateway/api/commission"
	"github.com/NpoolPlatform/inspire-gateway/api/reconciliation"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	inspire.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	inspire.RegisterGatewayServer(server, &Server{})
	archivement.Register(server)
	commission.Register(server)
	reconciliation.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := inspire.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := archivement.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := commission.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := reconciliation.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
