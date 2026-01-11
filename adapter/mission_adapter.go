package adapter

import (
	"github.com/hackathon-20260110/api/models"
	"gorm.io/gorm"
)

type MissionAdapter interface {
	GetMissionsByOwnerUserID(ownerUserID string) ([]models.Mission, error)
	GetMissionUnlocksByUserID(userID string) ([]models.MissionUnlock, error)
	GetMissionUnlock(missionID string, userID string) (*models.MissionUnlock, error)
	CreateMissionUnlock(missionUnlock models.MissionUnlock) error
}

type missionAdapter struct {
	db *gorm.DB
}

func NewMissionAdapter(db *gorm.DB) MissionAdapter {
	return &missionAdapter{db: db}
}

func (a *missionAdapter) GetMissionsByOwnerUserID(ownerUserID string) ([]models.Mission, error) {
	var missions []models.Mission
	if err := a.db.Where("mission_owner_user_id = ?", ownerUserID).Find(&missions).Error; err != nil {
		return nil, err
	}
	return missions, nil
}

func (a *missionAdapter) GetMissionUnlocksByUserID(userID string) ([]models.MissionUnlock, error) {
	var missionUnlocks []models.MissionUnlock
	if err := a.db.Where("unlocked_user_id = ?", userID).Find(&missionUnlocks).Error; err != nil {
		return nil, err
	}
	return missionUnlocks, nil
}

func (a *missionAdapter) GetMissionUnlock(missionID string, userID string) (*models.MissionUnlock, error) {
	var missionUnlock models.MissionUnlock
	if err := a.db.Where("mission_id = ? AND unlocked_user_id = ?", missionID, userID).First(&missionUnlock).Error; err != nil {
		return nil, err
	}
	return &missionUnlock, nil
}

func (a *missionAdapter) CreateMissionUnlock(missionUnlock models.MissionUnlock) error {
	return a.db.Create(&missionUnlock).Error
}
