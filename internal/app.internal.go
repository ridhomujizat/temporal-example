package internal

import (
	"context"
	"onx-outgoing-go/config"
	"onx-outgoing-go/internal/pkg/logger"
	"onx-outgoing-go/internal/pkg/middleware"
	"onx-outgoing-go/internal/pkg/rabbitmq"
	"onx-outgoing-go/internal/pkg/redis"
	"onx-outgoing-go/internal/repository"
	accountRepository "onx-outgoing-go/internal/repository/account"
	botService "onx-outgoing-go/internal/service/bot"
	"sync"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

func Setup(env config.Config, engine *gin.Engine, ctx context.Context, wg *sync.WaitGroup, redis redis.IRedis, rabbitMq *rabbitmq.ConnectionManager, publisher *rabbitmq.Publisher, temporalClient client.Client) {
	InitMiddelware(engine)
	e := engine.Group(BasePath())

	repository := &repository.Repository{
		Account: accountRepository.NewService(ctx, redis, rabbitMq, publisher),
	}

	InitRoutes(e, ctx, wg, redis)
	InitBotWorkflow(ctx, redis, rabbitMq, repository, temporalClient)

}

func BasePath() string {
	return "/api"
}

func InitMiddelware(e *gin.Engine) {
	e.Use(middleware.CorsMiddleware())
}

func InitRoutes(e *gin.RouterGroup, ctx context.Context, wg *sync.WaitGroup, redis redis.IRedis) {

	// /* init handler */
	// ExampleHandler := exampleHandler.NewHandler(service)
	// SessionHandler := sessionHandler.NewHandler(service)

	// /* init routes */
	// SessionHandler.NewRoutes(e)
	// ExampleHandler.NewRoutes(e)

	// Health check
	e.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

}
func InitBotWorkflow(ctx context.Context, redis redis.IRedis, rabbitMq *rabbitmq.ConnectionManager, repository *repository.Repository, temporal client.Client) {

	Botservice, err := botService.NewService(ctx, redis, rabbitMq, repository, temporal)
	if err != nil {
		logger.Error.Println("Failed to init service whatsapp", err)
	}
	if err := Botservice.Init(); err != nil {
		logger.Error.Println("Failed to init service whatsapp", err)
	}
	if err := Botservice.SubscribeWhatsappBot(); err != nil {
		logger.Error.Println("Failed to subcribe whatsapp service: ", err)
	}

}
