package models

import (
	"encoding/json"
	"time"
)

type ActionType string

const (
	ActionCreate  ActionType = "CREATE"
	ActionUpdate  ActionType = "UPDATE"
	ActionDelete  ActionType = "DELETE"
	ActionPayment ActionType = "PAYMENT"
)

type HistoryLog struct {
	ID            string          `gorm:"primaryKey" json:"id"`
	EntityID      string          `json:"entityId"`
	EntityType    string          `json:"entityType"`
	Action        ActionType      `json:"action"`
	Details       string          `json:"details"`
	Changes       json.RawMessage `gorm:"type:jsonb" json:"changes" swaggertype:"string"`
	PreviousValue json.RawMessage `gorm:"type:jsonb" json:"previousValue" swaggertype:"string"`
	NewValue      json.RawMessage `gorm:"type:jsonb" json:"newValue" swaggertype:"string"`
	Timestamp     time.Time       `json:"timestamp"`
	UserID        string          `json:"userId"`
	User          User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
