package models

import "time"

type Matching struct {
	ID string `gorm:"primaryKey" json:"id"`
	// User1ID < User2ID
	User1ID   string    `json:"user1_id" gorm:"not null"`
	User2ID   string    `json:"user2_id" gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
