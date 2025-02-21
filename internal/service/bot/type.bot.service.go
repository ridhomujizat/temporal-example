package botService

import (
	"context"
	"fmt"
	"onx-outgoing-go/internal/pkg/helper"
	"onx-outgoing-go/internal/pkg/rabbitmq"
	"onx-outgoing-go/internal/pkg/redis"
	"onx-outgoing-go/internal/repository"

	"go.temporal.io/sdk/client"
)

type Service struct {
	ctx           context.Context
	redis         redis.IRedis
	rabbitmq      *rabbitmq.ConnectionManager
	repository    *repository.Repository
	temporal      client.Client
	taskQueueName string
}

type IService interface {
	SubscribeWhatsappBot() error
	Init() error
}

func NewService(ctx context.Context, redis redis.IRedis, manager *rabbitmq.ConnectionManager, repository *repository.Repository, temp client.Client) (IService, error) {
	return &Service{
		ctx:           ctx,
		redis:         redis,
		rabbitmq:      manager,
		repository:    repository,
		temporal:      temp,
		taskQueueName: fmt.Sprintf("%s:bot-engine", helper.GetEnv("APP_TENANT")),
	}, nil
}
