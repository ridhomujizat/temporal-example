package model

import (
	"time"
)

type BotAccount struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	BotId       string    `gorm:"size:100;uniqueIndexx" json:"username"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   time.Time `gorm:"index" json:"deleted_at"`
}

type BotWorkflow struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      uint       `gorm:"not null" json:"name"`
	ParentId  *uint      `gorm:"default:null" json:"parent_id"`
	BotID     uint       `gorm:"not null" json:"bot_id"`
	Nodes     string     `gorm:"type:jsonb" json:"nodes"`
	Edges     string     `gorm:"type:jsonb" json:"edges"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time  `gorm:"index" json:"deleted_at"`
	Bot       BotAccount `gorm:"foreignKey:BotID" json:"bot"`
}
