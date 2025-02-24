package botRepository

import (
	"context"
	"onx-outgoing-go/internal/common/model"
	"onx-outgoing-go/internal/pkg/postgre"
	"onx-outgoing-go/internal/pkg/redis"
)

type Repository struct {
	ctx   context.Context
	redis redis.IRedis
	db    postgre.IPostgre
}

type IRepository interface {
	GetBot(idBot uint) (model.BotAccount, error)
	GetBotWorkflow(idBot uint) ([]model.BotWorkflow, error)
	UpdateBot(bot model.BotAccount) error
	CreateBot(bot model.BotAccount) error
	DeleteBot(bot model.BotAccount) error
	GetBotWorkflowById(idBot uint, idWorkflow string) (model.BotWorkflow, error)
	UpdateBotWorkflow(workflow model.BotWorkflow) error
	CreateBotWorkflow(workflow model.BotWorkflow) error
	DeleteBotWorkflow(workflow model.BotWorkflow) error
}

func NewService(ctx context.Context, redis redis.IRedis, db postgre.IPostgre) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}
