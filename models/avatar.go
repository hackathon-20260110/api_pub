package models

import "time"

type Avatar struct {
	ID                string    `gorm:"primaryKey" json:"id"`
	UserID            string    `json:"user_id" gorm:"not null"`
	AvatarIconURL     string    `json:"avatar_icon_url" gorm:"not null"`
	Prompt            string    `json:"prompt" gorm:"not null"`
	PersonalityTraits string    `json:"personality_traits" gorm:"type:jsonb"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
