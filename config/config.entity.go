package config

import (
	"context"
	"onx-outgoing-go/internal/common/enum"
	"onx-outgoing-go/internal/pkg/rabbitmq"
	"onx-outgoing-go/internal/pkg/redis"
	"sync"

	"go.temporal.io/sdk/client"
)

type Config struct {
	AppEnv        enum.EnvEnum `env:"APP_ENV" envDefault:"development"`
	AppPort       int          `env:"APP_PORT" envDefault:"8080"`
	AppTenant     string       `env:"APP_TENANT" envDefault:"tenant"`
	RedisHost     string       `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     int          `env:"REDIS_PORT" envDefault:"6379"`
	RedisPass     string       `env:"REDIS_PASS" envDefault:""`
	RedisPoolSize int          `env:"REDIS_POOL_SIZE" envDefault:"10"`
	RabbitHost    string       `env:"RABBIT_HOST" envDefault:"localhost"`
	RabbitPort    int          `env:"RABBIT_PORT" envDefault:"5672"`
	RabbitUser    string       `env:"RABBIT_USER" envDefault:"guest"`
	RabbitPass    string       `env:"RABBIT_PASS" envDefault:""`
}

type SetupServerDto struct {
	Ctx    *context.Context
	Cancel context.CancelFunc
	Wg     *sync.WaitGroup
	Env    *Config
	Rds    redis.IRedis
	Rb     *rabbitmq.ConnectionManager
	Pb     *rabbitmq.Publisher
	Tp     client.Client
}
