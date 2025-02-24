package model

import (
	"onx-outgoing-go/internal/common/enum"
	"time"
)

type AccountSetting struct {
	ID              uint                 `gorm:"primaryKey" json:"id"`
	Alias           string               `gorm:"type:varchar;column:alias" json:"alias"`
	Account         string               `gorm:"type:varchar;column:account" json:"account"`
	Key             string               `gorm:"type:text;column:key" json:"key"`
	ChannelId       enum.ChannelId       `gorm:"column:channel_id" json:"channel_id"`
	ChannelPlatform enum.ChannelPlatform `gorm:"type:varchar;column:channel_platform" json:"channel_platform"`
	BaseURL         string               `gorm:"type:text;column:base_url" json:"base_url"`
	BotID           uint                 `gorm:"column:bot_id" json:"bot_id"`
	CreatedAt       time.Time            `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time            `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       time.Time            `gorm:"index" json:"deleted_at"`
	Bot             BotAccount           `gorm:"foreignKey:BotID" json:"bot"`
}
