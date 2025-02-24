package accountRepository

import (
	"context"
	"onx-outgoing-go/internal/pkg/postgre"
	"onx-outgoing-go/internal/pkg/redis"
)

type Repository struct {
	ctx   context.Context
	redis redis.IRedis
	db    postgre.IPostgre
}

type IService interface {
}

func NewService(ctx context.Context, redis redis.IRedis, db postgre.IPostgre) IService {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}
