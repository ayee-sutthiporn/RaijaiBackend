package models

import (
	"encoding/json"
	"time"
)

type DebtType string

const (
	DebtTypeLent     DebtType = "LENT"
	DebtTypeBorrowed DebtType = "BORROWED"
)

type Debt struct {
	ID              string          `gorm:"primaryKey" json:"id"`
	Title           string          `json:"title"`
	TotalAmount     float64         `json:"total_amount"`
	RemainingAmount float64         `json:"remaining_amount"`
	Type            DebtType        `json:"type"`
	PersonName      string          `json:"person_name"`
	DueDate         time.Time       `json:"due_date"`
	WalletID        *string         `json:"wallet_id"`
	Wallet          *Wallet         `gorm:"foreignKey:WalletID" json:"wallet"`
	IsInstallment   bool            `json:"is_installment"`
	InstallmentPlan json.RawMessage `gorm:"type:jsonb" json:"installment_plan"` // Using RawMessage for JSONB
	Remark          string          `json:"remark"`
	AutoDeduct      bool            `json:"auto_deduct"`
	CreatedAt       time.Time       `json:"created_at"`
}
