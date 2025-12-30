package models

import (
	"time"
)

type DebtType string

const (
	DebtTypeLent     DebtType = "LENT"
	DebtTypeBorrowed DebtType = "BORROWED"
)

type InstallmentPlan struct {
	TotalMonths   int       `json:"totalMonths"`
	PaidMonths    int       `json:"paidMonths"`
	InterestRate  float64   `json:"interestRate"`
	MonthlyAmount float64   `json:"monthlyAmount"`
	StartDate     time.Time `json:"startDate"`
}

type Debt struct {
	ID              string           `gorm:"primaryKey" json:"id"`
	Title           string           `json:"title"`
	TotalAmount     float64          `json:"totalAmount"`
	RemainingAmount float64          `json:"remainingAmount"`
	Type            DebtType         `json:"type"`
	PersonName      string           `json:"personName"`
	DueDate         time.Time        `json:"dueDate"`
	WalletID        *string          `json:"walletId"`
	Wallet          *Wallet          `gorm:"foreignKey:WalletID" json:"wallet"`
	IsInstallment   bool             `json:"isInstallment"`
	InstallmentPlan *InstallmentPlan `gorm:"serializer:json" json:"installmentPlan"`
	Remark          string           `json:"remark"`
	AutoDeduct      bool             `json:"autoDeduct"`
	CreatedAt       time.Time        `json:"createdAt"`
}
