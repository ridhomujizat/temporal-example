package postgre

import (
	"context"
	"fmt"
	"onx-outgoing-go/internal/common/model"
	"onx-outgoing-go/internal/pkg/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Setup(ctx context.Context, config *Config) (IPostgre, error) {
	clientCtx, cancel := context.WithCancel(ctx)

	p := &Client{
		cancel: cancel,
		ctx:    clientCtx,
		config: config,
	}

	if err := p.connect(); err != nil {
		cancel() // Ensure cleanup if initialization fails
		logger.Error.Println(err.Error())
		return nil, fmt.Errorf("failed to connect to postgre: %w", err)
	}

	// Start the reconnect handler
	go p.reconnectHandler()

	return p, nil
}

func (p *Client) connect() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		p.config.Host,
		p.config.Username,
		p.config.Password,
		p.config.Name,
		p.config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logger.Error.Println("Error connecting to database:", err)
		return fmt.Errorf("error connecting to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error.Println("Error getting database instance:", err)
		return fmt.Errorf("error getting database instance: %w", err)
	}

	logger.Info.Println("Connected to database")
	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	p.db = db

	// err = p.AutoMigrateEnum()
	// if err != nil {
	// 	logger.Warning.Println("Error migrating enum database:", err)
	// }
	err = p.Migrate()
	if err != nil {
		logger.Error.Println("Error migrating database:", err)
	}

	return nil
}

func (p *Client) Close() error {
	p.cancel()
	if sqlDB, err := p.db.DB(); err == nil {
		return sqlDB.Close()
	}
	return nil
}

func (p *Client) Ping() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetDB returns the underlying gorm.DB instance
func (p *Client) GetDB() *gorm.DB {
	return p.db
}

func (p *Client) reconnectHandler() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			if sqlDB, err := p.db.DB(); err != nil || sqlDB.Ping() != nil {
				logger.Error.Println("Database connection lost, attempting to reconnect...")
				if err := p.connect(); err != nil {
					logger.Error.Println("Reconnection failed:", err)
				} else {
					logger.Info.Println("Successfully reconnected to database")
				}
			}
		}
	}
}

func (p *Client) AutoMigrateEnum() error {
	// bot_platform_enum := "'botpress', 'typebot'"
	// channel_source_enum := "'whatsapp', 'fbmessenger', 'livechat', 'igdm', 'telegram'"
	// channel_platform_enum := "'socioconnect', 'maytapi', 'octopushchat', 'official'"
	// channel_id_enum := "'12', '3', '7', '16', '5'"
	// omnichannel_enum := "'onx', 'on5', 'on4'"
	// account_platform_enum := "'whatsapp_socio', 'fbm_socio', 'whatsapp_maytapi', 'botpress', 'livechat_octopushchat', 'igdm_socio', 'official'"
	// account_type_enum := "'igdm', 'octopushchat', 'fbmessenger', 'whatsapp', 'botpress', 'telegram'"

	// if err := p.db.Exec(fmt.Sprintf(`
	// 	CREATE SEQUENCE IF NOT EXISTS id_seq;
	// 	DROP TYPE IF EXISTS "public"."bot_platform_enum";
	// 	CREATE TYPE "public"."bot_platform_enum" AS ENUM (%s);
	// 	DROP TYPE IF EXISTS "public"."channel_source_enum";
	// 	CREATE TYPE "public"."channel_source_enum" AS ENUM (%s);
	// 	DROP TYPE IF EXISTS "public"."channel_platform_enum";
	// 	CREATE TYPE "public"."channel_platform_enum" AS ENUM (%s);
	// 	DROP TYPE IF EXISTS "public"."channel_id_enum";
	// 	CREATE TYPE "public"."channel_id_enum" AS ENUM (%s);
	// 	DROP TYPE IF EXISTS "public"."omnichannel_enum";
	// 	CREATE TYPE "public"."omnichannel_enum" AS ENUM (%s);
	// 	DROP TYPE IF EXISTS "public"."account_platform_enum";
	// 	CREATE TYPE "public"."account_platform_enum" AS ENUM (%s);
	// 	DROP TYPE IF EXISTS "public"."account_type_enum";
	// 	CREATE TYPE "public"."account_type_enum" AS ENUM (%s);
	// `,
	// 	bot_platform_enum,
	// 	channel_source_enum,
	// 	channel_platform_enum,
	// 	channel_id_enum,
	// 	omnichannel_enum,
	// 	account_platform_enum,
	// 	account_type_enum)).Error; err != nil {
	// 	return fmt.Errorf("error migrating database: %w", err)
	// }
	return nil
}

func (p *Client) Migrate() error {
	err := p.db.AutoMigrate(
		&model.BotAccount{},
		// &model.FlowSetting{},
		// &model.Session{},
		// &model.SessionHistory{},
		// &model.InteractionHistory{},
		// &model.Report{},
		// &model.ReportCsat{},
		// &model.ConfigKeyIncoming{},
		// &model.Template{},
	)
	if err != nil {
		return fmt.Errorf("error migrating database: %w", err)

	}
	return nil
}
