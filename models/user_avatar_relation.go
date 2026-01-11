package models

import "time"

type UserAvatarRelation struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	UserID        string    `json:"user_id" gorm:"not null"`
	AvatarID      string    `json:"avatar_id" gorm:"not null"`
	MatchingPoint int       `json:"matching_point" gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
