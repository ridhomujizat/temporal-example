package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	ID        uint                  `gorm:"primaryKey" json:"id"`
	Name      string                `gorm:"not null" json:"name"`
	ParentId  *uint                 `gorm:"default:null" json:"parent_id"`
	BotID     uint                  `gorm:"not null" json:"bot_id"`
	Nodes     *BotWorkflowNode      `gorm:"type:jsonb" json:"nodes"`
	Edges     BotWorkflowEdgesSlice `gorm:"type:jsonb" json:"edges"`
	CreatedAt time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time             `gorm:"index" json:"deleted_at"`
	Bot       BotAccount            `gorm:"foreignKey:BotID" json:"bot"`
}

type BotWorkflowNode struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Nodes       []Node `json:"nodes"`
	Version     string `json:"version"`
	Description string `json:"description"`
	StartNodeID string `json:"startNodeId"`
}

// Implement the driver.Valuer interface for BotWorkflowNode
func (n BotWorkflowNode) Value() (driver.Value, error) {
	return json.Marshal(n)
}

// Implement the sql.Scanner interface for *BotWorkflowNode
func (n *BotWorkflowNode) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal BotWorkflowNode: %v", value)
	}
	return json.Unmarshal(bytes, n)
}

// Create a custom type for a slice of BotWorkflowEdges
type BotWorkflowEdgesSlice []BotWorkflowEdges

// Implement the driver.Valuer interface for BotWorkflowEdgesSlice
func (bes BotWorkflowEdgesSlice) Value() (driver.Value, error) {
	return json.Marshal(bes)
}

// Implement the sql.Scanner interface for BotWorkflowEdgesSlice
func (bes *BotWorkflowEdgesSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal BotWorkflowEdges value: %v", value)
	}
	return json.Unmarshal(bytes, bes)
}

type Node struct {
	ID               string           `json:"id"`
	Title            string           `json:"title"`
	Blocks           []Block          `json:"blocks"`
	GraphCoordinates GraphCoordinates `json:"graphCoordinates"`
}

type Block struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Content    string   `json:"content"`
	Choices    []Choice `json:"choices,omitempty"`
	IsDropdown *bool    `json:"isDropdown,omitempty"`
	NextEdgeID *string  `json:"nextEdgeId,omitempty"`
}

type Choice struct {
	Value      string `json:"value"`
	Content    string `json:"content"`
	NextEdgeID string `json:"nextEdgeId"`
}

type GraphCoordinates struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type BotWorkflowEdges struct {
	ID   string `json:"id"`
	To   To     `json:"to"`
	From From   `json:"from"`
}

type From struct {
	ItemID string `json:"itemId"`
	NodeID string `json:"nodeId"`
}

type To struct {
	Type   string `json:"type"`
	NodeID string `json:"nodeId"`
}
