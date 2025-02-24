package model

import (
	"time"
)

type BotAccount struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	BotId       string    `gorm:"size:100;uniqueIndexx" json:"username"`
	ApiKey      string    `gorm:"size:255" json:"api_key"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   time.Time `gorm:"index" json:"deleted_at"`
}

func (BotAccount) TableName() string {
	return "bot_accounts"
}
