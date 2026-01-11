package adapter

import (
	"context"
	"fmt"

	"github.com/hackathon-20260110/api/driver"
	"github.com/hackathon-20260110/api/models"
	"gorm.io/gorm"
)

type ProfileAdapter interface {
	// UserInfo関連
	CreateUserInfo(userInfo models.UserInfo) (models.UserInfo, error)
	UpdateUserInfo(id string, userInfo models.UserInfo) (models.UserInfo, error)
	DeleteUserInfo(id string) error
	GetUserInfoByUserID(userID string) ([]models.UserInfo, error)
	GetUserInfoByID(id string) (models.UserInfo, error)

	// Mission関連
	CreateMission(mission models.Mission) (models.Mission, error)
	UpdateMission(id string, mission models.Mission) (models.Mission, error)
	DeleteMissionByUserInfoID(userInfoID string) error
	GetMissionsByOwnerID(ownerID string) ([]models.Mission, error)
	GetMissionByUserInfoID(userInfoID string) (models.Mission, error)

	// MissionUnlock関連
	CreateMissionUnlock(unlock models.MissionUnlock) (models.MissionUnlock, error)
	GetMissionUnlocksByUserID(userID string, missionIDs []string) ([]models.MissionUnlock, error)
	CheckMissionUnlocked(missionID string, userID string) (bool, error)

	// マッチングポイント取得（PostgreSQLから）
	GetMatchingScore(ctx context.Context, userID string, partnerUserID string) (int, error)
}

func NewProfileAdapter() ProfileAdapter {
	db := driver.NewPsql()
	return &profileAdapter{db: db}
}

type profileAdapter struct {
	db *gorm.DB
}

// CreateUserInfo UserInfo作成
func (a *profileAdapter) CreateUserInfo(userInfo models.UserInfo) (models.UserInfo, error) {
	if err := a.db.Create(&userInfo).Error; err != nil {
		return models.UserInfo{}, err
	}
	return userInfo, nil
}

// UpdateUserInfo UserInfo更新
func (a *profileAdapter) UpdateUserInfo(id string, userInfo models.UserInfo) (models.UserInfo, error) {
	if err := a.db.Model(&models.UserInfo{}).Where("id = ?", id).Updates(userInfo).Error; err != nil {
		return models.UserInfo{}, err
	}
	var updated models.UserInfo
	if err := a.db.Where("id = ?", id).First(&updated).Error; err != nil {
		return models.UserInfo{}, err
	}
	return updated, nil
}

// DeleteUserInfo UserInfo削除（関連するMissionも削除）
func (a *profileAdapter) DeleteUserInfo(id string) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		// 関連するMissionを削除
		if err := tx.Where("user_info_id = ?", id).Delete(&models.Mission{}).Error; err != nil {
			return err
		}
		// UserInfoを削除
		if err := tx.Where("id = ?", id).Delete(&models.UserInfo{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetUserInfoByUserID ユーザーIDでUserInfo一覧を取得
func (a *profileAdapter) GetUserInfoByUserID(userID string) ([]models.UserInfo, error) {
	var userInfoList []models.UserInfo
	if err := a.db.Where("user_id = ?", userID).Find(&userInfoList).Error; err != nil {
		return nil, err
	}
	return userInfoList, nil
}

// GetUserInfoByID IDでUserInfoを取得
func (a *profileAdapter) GetUserInfoByID(id string) (models.UserInfo, error) {
	var userInfo models.UserInfo
	if err := a.db.Where("id = ?", id).First(&userInfo).Error; err != nil {
		return models.UserInfo{}, err
	}
	return userInfo, nil
}

// CreateMission Mission作成
func (a *profileAdapter) CreateMission(mission models.Mission) (models.Mission, error) {
	if err := a.db.Create(&mission).Error; err != nil {
		return models.Mission{}, err
	}
	return mission, nil
}

// UpdateMission Mission更新
func (a *profileAdapter) UpdateMission(id string, mission models.Mission) (models.Mission, error) {
	if err := a.db.Model(&models.Mission{}).Where("id = ?", id).Updates(mission).Error; err != nil {
		return models.Mission{}, err
	}
	var updated models.Mission
	if err := a.db.Where("id = ?", id).First(&updated).Error; err != nil {
		return models.Mission{}, err
	}
	return updated, nil
}

// DeleteMissionByUserInfoID UserInfoIDでMissionを削除
func (a *profileAdapter) DeleteMissionByUserInfoID(userInfoID string) error {
	return a.db.Where("user_info_id = ?", userInfoID).Delete(&models.Mission{}).Error
}

// GetMissionsByOwnerID オーナーIDでMission一覧を取得
func (a *profileAdapter) GetMissionsByOwnerID(ownerID string) ([]models.Mission, error) {
	var missions []models.Mission
	if err := a.db.Where("mission_owner_user_id = ?", ownerID).Find(&missions).Error; err != nil {
		return nil, err
	}
	return missions, nil
}

// GetMissionByUserInfoID UserInfoIDでMissionを取得
func (a *profileAdapter) GetMissionByUserInfoID(userInfoID string) (models.Mission, error) {
	var mission models.Mission
	if err := a.db.Where("user_info_id = ?", userInfoID).First(&mission).Error; err != nil {
		return models.Mission{}, err
	}
	return mission, nil
}

// CreateMissionUnlock MissionUnlock作成
func (a *profileAdapter) CreateMissionUnlock(unlock models.MissionUnlock) (models.MissionUnlock, error) {
	if err := a.db.Create(&unlock).Error; err != nil {
		return models.MissionUnlock{}, err
	}
	return unlock, nil
}

// GetMissionUnlocksByUserID ユーザーIDとMissionIDリストでMissionUnlock一覧を取得
func (a *profileAdapter) GetMissionUnlocksByUserID(userID string, missionIDs []string) ([]models.MissionUnlock, error) {
	if len(missionIDs) == 0 {
		return []models.MissionUnlock{}, nil
	}
	var unlocks []models.MissionUnlock
	if err := a.db.Where("unlocked_user_id = ? AND mission_id IN ?", userID, missionIDs).Find(&unlocks).Error; err != nil {
		return nil, err
	}
	return unlocks, nil
}

// CheckMissionUnlocked 特定のMissionが特定のユーザーによって解禁されているか確認
func (a *profileAdapter) CheckMissionUnlocked(missionID string, userID string) (bool, error) {
	var count int64
	if err := a.db.Model(&models.MissionUnlock{}).
		Where("mission_id = ? AND unlocked_user_id = ?", missionID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetMatchingScore マッチングポイントを取得（PostgreSQLから）
// 相手ユーザーが所有するアバターとのUserAvatarRelationのマッチングポイントを合計して返す
func (a *profileAdapter) GetMatchingScore(ctx context.Context, userID string, partnerUserID string) (int, error) {
	var totalScore int

	err := a.db.Model(&models.UserAvatarRelation{}).
		Joins("INNER JOIN avatars ON user_avatar_relations.avatar_id = avatars.id").
		Where("user_avatar_relations.user_id = ? AND avatars.user_id = ?", userID, partnerUserID).
		Select("COALESCE(SUM(user_avatar_relations.matching_point), 0) as total_score").
		Scan(&totalScore).Error

	if err != nil {
		return 0, fmt.Errorf("failed to get matching score: %w", err)
	}

	return totalScore, nil
}
