package main

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/action"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/inspire-gateway/api"
	"github.com/NpoolPlatform/inspire-middleware/pkg/db"

	apicli "github.com/NpoolPlatform/basal-middleware/pkg/client/api"
	migrator "github.com/NpoolPlatform/inspire-gateway/pkg/migrator"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	cli "github.com/urfave/cli/v2"

	"google.golang.org/grpc"

	"github.com/NpoolPlatform/inspire-gateway/pkg/pubsub"
)

var runCmd = &cli.Command{
	Name:    "run",
	Aliases: []string{"s"},
	Usage:   "Run the daemon",
	Action: func(c *cli.Context) error {
		return action.Run(
			c.Context,
			run,
			rpcRegister,
			rpcGatewayRegister,
			watch,
		)
	},
}

func run(ctx context.Context) error {
	if err := db.Init(); err != nil {
		return err
	}
	if err := migrator.Migrate(ctx); err != nil {
		return err
	}
	return pubsub.Subscribe(ctx)
}

func shutdown(ctx context.Context) {
	<-ctx.Done()
	logger.Sugar().Infow(
		"Watch",
		"State", "Done",
		"Error", ctx.Err(),
	)
	_ = pubsub.Shutdown(ctx)
}

func watch(ctx context.Context, cancel context.CancelFunc) error {
	go shutdown(ctx)
	return nil
}

func rpcRegister(server grpc.ServiceRegistrar) error {
	api.Register(server)

	apicli.RegisterGRPC(server)

	return nil
}

func rpcGatewayRegister(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	err := api.RegisterGateway(mux, endpoint, opts)
	if err != nil {
		return err
	}

	_ = apicli.Register(mux)

	return nil
}
