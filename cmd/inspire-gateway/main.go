package main

import (
	"fmt"
	"os"

	servicename "github.com/NpoolPlatform/inspire-gateway/pkg/servicename"

	"github.com/NpoolPlatform/go-service-framework/pkg/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	mysqlconst "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	rabbitmqconst "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/const"
	redisconst "github.com/NpoolPlatform/go-service-framework/pkg/redis/const"

	inspireconst "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/message/const"
	archivementmgrconst "github.com/NpoolPlatform/inspire-manager/pkg/message/const"

	cli "github.com/urfave/cli/v2"
)

func main() {
	commands := cli.Commands{
		runCmd,
	}

	description := fmt.Sprintf("my %v service cli\nFor help on any individual command run <%v COMMAND -h>\n",
		servicename.ServiceName, servicename.ServiceName)
	err := app.Init(
		servicename.ServiceName,
		description,
		"",
		"",
		"./",
		nil,
		commands,
		mysqlconst.MysqlServiceName,
		rabbitmqconst.RabbitMQServiceName,
		redisconst.RedisServiceName,
		inspireconst.ServiceName,
		archivementmgrconst.ServiceName,
	)
	if err != nil {
		logger.Sugar().Errorf("fail to create %v: %v", servicename.ServiceName, err)
		return
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Sugar().Errorf("fail to run %v: %v", servicename.ServiceName, err)
	}
}
