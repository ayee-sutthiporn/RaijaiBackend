package models

import (
	"time"
)

type Book struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `gorm:"index" json:"ownerId"`
	Owner       User      `gorm:"foreignKey:OwnerID" json:"owner"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type BookMember struct {
	BookID    string    `gorm:"primaryKey" json:"bookId"`
	UserID    string    `gorm:"primaryKey" json:"userId"`
	Role      string    `json:"role"` // OWNER, EDITOR, VIEWER
	JoinedAt  time.Time `json:"joinedAt"`
	Book      Book      `gorm:"foreignKey:BookID" json:"book"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}
