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
	WalletID    string          `json:"walletId"`
	Wallet      Wallet          `gorm:"foreignKey:WalletID" json:"wallet"`
	ToWalletID  *string         `json:"toWalletId"`
	ToWallet    *Wallet         `gorm:"foreignKey:ToWalletID" json:"toWallet"`
	Amount      float64         `json:"amount"`
	Type        TransactionType `json:"type"`
	CategoryID  string          `json:"categoryId"`
	Category    *Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Description string          `json:"description"`
	Date        DateOnly        `gorm:"type:date" json:"date"`
	CreatedByID string          `gorm:"index" json:"createdById"`
	CreatedBy   User            `gorm:"foreignKey:CreatedByID" json:"createdBy"`
	CreatedAt   time.Time       `json:"createdAt"`
}
