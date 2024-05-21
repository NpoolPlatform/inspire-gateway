package api

import (
	"context"

	inspire "github.com/NpoolPlatform/message/npool/inspire/gw/v1"

	"github.com/NpoolPlatform/inspire-gateway/api/achievement"
	appcommissionconfig "github.com/NpoolPlatform/inspire-gateway/api/app/commission/config"
	appconfig "github.com/NpoolPlatform/inspire-gateway/api/app/config"
	appgoodcommissionconfig "github.com/NpoolPlatform/inspire-gateway/api/app/good/commission/config"
	coinallocated "github.com/NpoolPlatform/inspire-gateway/api/coin/allocated"
	coinconfig "github.com/NpoolPlatform/inspire-gateway/api/coin/config"
	"github.com/NpoolPlatform/inspire-gateway/api/commission"
	"github.com/NpoolPlatform/inspire-gateway/api/coupon"
	"github.com/NpoolPlatform/inspire-gateway/api/coupon/allocated"
	cashcontrol "github.com/NpoolPlatform/inspire-gateway/api/coupon/app/cashcontrol"
	appgoodscope "github.com/NpoolPlatform/inspire-gateway/api/coupon/app/scope"
	"github.com/NpoolPlatform/inspire-gateway/api/coupon/scope"
	"github.com/NpoolPlatform/inspire-gateway/api/event"
	"github.com/NpoolPlatform/inspire-gateway/api/invitation/invitationcode"
	"github.com/NpoolPlatform/inspire-gateway/api/invitation/registration"
	"github.com/NpoolPlatform/inspire-gateway/api/reconcile"
	taskconfig "github.com/NpoolPlatform/inspire-gateway/api/task/config"
	usercoinreward "github.com/NpoolPlatform/inspire-gateway/api/user/coin/reward"
	userreward "github.com/NpoolPlatform/inspire-gateway/api/user/reward"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	inspire.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	inspire.RegisterGatewayServer(server, &Server{})
	achievement.Register(server)
	commission.Register(server)
	reconcile.Register(server)
	coupon.Register(server)
	allocated.Register(server)
	scope.Register(server)
	appgoodscope.Register(server)
	invitationcode.Register(server)
	registration.Register(server)
	event.Register(server)
	cashcontrol.Register(server)
	appgoodcommissionconfig.Register(server)
	appcommissionconfig.Register(server)
	appconfig.Register(server)
	taskconfig.Register(server)
	coinconfig.Register(server)
	coinallocated.Register(server)
	userreward.Register(server)
	usercoinreward.Register(server)
}

//nolint:gocyclo
func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := inspire.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := achievement.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := commission.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := reconcile.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := coupon.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appgoodscope.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := allocated.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := scope.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := invitationcode.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := registration.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := event.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := cashcontrol.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appgoodcommissionconfig.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appcommissionconfig.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appconfig.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := taskconfig.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := coinconfig.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := coinallocated.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := userreward.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := usercoinreward.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
