package models

import "time"

// PasswordResetToken stores a single-use, time-limited token used to
// authorize a password reset. Only the SHA-256 hash of the raw token is
// persisted (see utils.HashToken) so a DB leak does not expose usable tokens.
type PasswordResetToken struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"index;not null" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	TokenHash string    `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `json:"-"`
	Used      bool      `gorm:"default:false" json:"-"`
	CreatedAt time.Time `json:"-"`
}
