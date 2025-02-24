package internal

import (
	"context"
	"onx-outgoing-go/config"
	"onx-outgoing-go/internal/pkg/logger"
	"onx-outgoing-go/internal/pkg/postgre"
	"onx-outgoing-go/internal/pkg/rabbitmq"
	"onx-outgoing-go/internal/pkg/redis"
	"onx-outgoing-go/internal/repository"
	botService "onx-outgoing-go/internal/service/bot"
	"sync"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

func Setup(env config.Config, engine *gin.Engine, ctx context.Context, wg *sync.WaitGroup, redis redis.IRedis, rabbitMq *rabbitmq.ConnectionManager, publisher *rabbitmq.Publisher, temporalClient client.Client, db postgre.IPostgre) {

	repository := &repository.Repository{
		// Account: accountRepository.NewService(ctx, redis, rabbitMq, publisher),
	}

	InitServicesOmnix(ctx, redis, rabbitMq, repository, temporalClient)

}

func InitServicesOmnix(ctx context.Context, redis redis.IRedis, rabbitMq *rabbitmq.ConnectionManager, repository *repository.Repository, temporal client.Client) {
	// err := temporal.InitTemporalConnection(ctx)
	// if err != nil {
	// 	fmt.Println("Unable to create Temporal client", err)
	// 	return
	// }

	whatsappp, err := botService.NewService(ctx, redis, rabbitMq, repository, temporal)
	if err != nil {
		logger.Error.Println("Failed to init service whatsapp", err)
	}
	if err := whatsappp.Init(); err != nil {
		logger.Error.Println("Failed to init service whatsapp", err)
	}
	if err := whatsappp.SubscribeWhatsappBot(); err != nil {
		logger.Error.Println("Failed to subcribe whatsapp service: ", err)
	}

	// defer temporal.CloseTemporalConnection()

}
