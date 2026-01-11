package models

import "time"

type User struct {
	ID                    string    `json:"id" gorm:"primaryKey"` // firebaseUIDを主キーにする
	DisplayName           string    `json:"display_name" gorm:"not null"`
	Gender                string    `json:"gender" gorm:"not null"` // male, female, other
	BirthDate             time.Time `json:"birth_date" gorm:"not null"`
	Bio                   string    `json:"bio" gorm:"not null"`
	ProfileImageURL       string    `json:"profile_image_url" gorm:"not null"`
	IsOnboardingCompleted bool      `json:"is_onboarding_completed" gorm:"default:false"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
