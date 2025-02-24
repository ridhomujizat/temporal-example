package postgre

import (
	"context"

	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Username string
	Password string
	Name     string
	Port     string
}

type Client struct {
	db     *gorm.DB
	cancel context.CancelFunc
	ctx    context.Context
	config *Config
}

type IPostgre interface {
	Close() error
	Ping() error
	GetDB() *gorm.DB
}
