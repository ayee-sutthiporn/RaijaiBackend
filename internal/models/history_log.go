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
	EntityID      string          `json:"entity_id"`
	EntityType    string          `json:"entity_type"`
	Action        ActionType      `json:"action"`
	Details       string          `json:"details"`
	Changes       json.RawMessage `gorm:"type:jsonb" json:"changes" swaggertype:"string"`
	PreviousValue json.RawMessage `gorm:"type:jsonb" json:"previous_value" swaggertype:"string"`
	NewValue      json.RawMessage `gorm:"type:jsonb" json:"new_value" swaggertype:"string"`
	Timestamp     time.Time       `json:"timestamp"`
}
