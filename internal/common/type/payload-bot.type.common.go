package types

import (
	"onx-outgoing-go/internal/common/enum"
	"onx-outgoing-go/internal/common/model"
)

type PayloadBot struct {
	CustMessage interface{} `json:"incoming"`
	Value       string      `json:"value"`
	MetaData    MetaData    `json:"metadata"`
}

// Incoming message from user ========================
type MetaData struct {
	// BotEndpoint     string               `json:"bot_endpoint" validate:"omitempty,url"`
	// BotAccount      string               `json:"bot_account" validate:"omitempty"`
	// BotAlias        string               `json:"bot_alias" validate:"omitempty"`
	AccountId       string               `json:"accountId" validate:"required"`
	UniqueId        string               `json:"unique_id" validate:"required"`
	CustName        string               `json:"cust_name" validate:"omitempty"`
	TenantId        string               `json:"tenant_id" validate:"omitempty"`
	ChannelId       enum.ChannelId       `json:"channel_id" validate:"required,oneof=12 3 7 5"`
	ChannelPlatform enum.ChannelPlatform `json:"channel_platform" validate:"required,oneof=socioconnect maytapi octopushchat official"`
	ChannelSources  enum.ChannelSources  `json:"channel_sources" validate:"required,oneof=whatsapp fbmessenger livechat telegram"`
	DateTimestamp   string               `json:"date_timestamp" validate:"required"`
	// Sid             string               `json:"sid,omitempty" validate:"omitempty"`
	// NewSession      bool                 `json:"new_session,omitempty" validate:"omitempty"`
	CustMessage     string          `json:"cust_message,omitempty" validate:"omitempty"`
	IsLocation      bool            `json:"isLocation,omitempty" validate:"omitempty"`
	LocationPayload LocationPayload `json:"locationPayload,omitempty" validate:"omitempty"`
}

type LocationPayload struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

// History chat ========================
type ResultWorkflowBot struct {
	AccountId       string            `json:"accountId" validate:"required"`
	UniqueId        string            `json:"unique_id" validate:"required"`
	ResultBlockChat []ResultBlockChat `json:"result_block_chat" validate:"required"`
	Error           string            `json:"error" validate:"omitempty"`
}

type ResultBlockChat struct {
	ID          string           `json:"id" validate:"required"`
	Node        string           `json:"node" validate:"required"`
	NextId      model.To         `json:"next_id" validate:"required"`
	HistoryChat []HistoryChatBot `json:"history_chat" validate:"required"`
}

type HistoryChatBot struct {
	From          string `json:"from" validate:"required"`
	Type          string `json:"type" validate:"required"`
	Message       string `json:"message" validate:"required"`
	DateTimestamp string `json:"date_timestamp" validate:"required"`
}
