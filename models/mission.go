package models

import "time"

type Mission struct {
	ID                      string    `gorm:"primaryKey" json:"id"`
	MissionOwnerUserID      string    `json:"mission_owner_user_id" gorm:"not null"`
	UserInfoID              string    `json:"user_info_id" gorm:"not null"`
	ThresholdPointCondition *int      `json:"threshold_point_condition"`
	UnlockCondition         *string   `json:"unlock_condition"`
	CreatedAt               time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
