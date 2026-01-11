package models

import "time"

type MissionUnlock struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	MissionID      string    `json:"mission_id" gorm:"not null"`
	UnlockedUserID string    `json:"unlocked_user_id" gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
