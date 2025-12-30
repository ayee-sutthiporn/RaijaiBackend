package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "INCOME"
	CategoryTypeExpense CategoryType = "EXPENSE"
)

type Category struct {
	ID        string       `gorm:"primaryKey" json:"id"`
	Name      string       `json:"name"`
	Type      CategoryType `json:"type"`
	Color     string       `json:"color"`
	Icon      string       `json:"icon"`
	CreatedAt time.Time    `json:"createdAt"`
}

type WalletType string

const (
	WalletTypeCash       WalletType = "CASH"
	WalletTypeBank       WalletType = "BANK"
	WalletTypeCreditCard WalletType = "CREDIT_CARD"
)

type Wallet struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	Type      WalletType `json:"type"`
	Balance   float64    `json:"balance"`
	Currency  string     `json:"currency"`
	Color     string     `json:"color"`
	OwnerID   string     `json:"ownerId"`
	Owner     User       `gorm:"foreignKey:OwnerID" json:"owner"`
	CreatedAt time.Time  `json:"createdAt"`
}
