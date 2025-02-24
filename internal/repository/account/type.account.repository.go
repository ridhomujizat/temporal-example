package accountRepository

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
	GetAccountSetting(account string) (model.AccountSetting, error)
	UpdateAccountSetting(account string, setting model.AccountSetting) error
	CreateAccountSetting(setting model.AccountSetting) error
	DeleteAccountSetting(account string) error
	GetBotByAccount(account string) (*model.BotAccount, error)
}

func NewService(ctx context.Context, redis redis.IRedis, db postgre.IPostgre) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}
