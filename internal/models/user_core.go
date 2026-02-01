package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex" json:"username"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"` // Don't expose password
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`   // ADMIN, USER
	Status    string    `json:"status"` // ACTIVE, INACTIVE
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserRegistration struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
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
	UserID    string       `gorm:"index" json:"userId"`
	User      User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
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
	OwnerID   string     `gorm:"index" json:"ownerId"`
	Owner     User       `gorm:"foreignKey:OwnerID" json:"owner"`
	CreatedAt time.Time  `json:"createdAt"`
}
