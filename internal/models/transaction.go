package models

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "INCOME"
	TransactionTypeExpense  TransactionType = "EXPENSE"
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

type Transaction struct {
	ID          string          `gorm:"primaryKey" json:"id"`
	WalletID    string          `json:"wallet_id"`
	Wallet      Wallet          `gorm:"foreignKey:WalletID" json:"wallet"`
	ToWalletID  *string         `json:"to_wallet_id"` // Pointer for nullable
	ToWallet    *Wallet         `gorm:"foreignKey:ToWalletID" json:"to_wallet"`
	Amount      float64         `json:"amount"`
	Type        TransactionType `json:"type"`
	Category    string          `json:"category"` // Assuming Name or ID stored as string
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	CreatedByID string          `json:"created_by_id"`
	CreatedBy   User            `gorm:"foreignKey:CreatedByID" json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
}
