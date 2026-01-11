package adapter

import (
	"github.com/hackathon-20260110/api/driver"
	"github.com/hackathon-20260110/api/models"
	"gorm.io/gorm"
)

type DiagnosisAdapter interface {
	// 診断履歴の保存（ポインタ渡しでGORMの自動設定値を反映）
	CreateDiagnosisHistory(history *models.DiagnosisHistory) error

	// ユーザーの診断履歴取得
	GetDiagnosisHistoryByUserID(userID string) ([]models.DiagnosisHistory, error)

	// 特定の診断履歴取得
	GetDiagnosisHistoryByID(id string) (models.DiagnosisHistory, error)

	// UserAvatarRelationのポイント更新
	UpdateUserAvatarRelationPoints(userID, avatarID string, additionalPoints int) error

	// UserAvatarRelationの取得（存在チェック用）
	GetUserAvatarRelation(userID, avatarID string) (models.UserAvatarRelation, error)

	// UserAvatarRelationの作成
	CreateUserAvatarRelation(relation models.UserAvatarRelation) error

	// Avatarの存在確認
	GetAvatarByID(id string) (models.Avatar, error)

	GetByUserID(userID string) (models.Avatar, error)
}

func NewDiagnosisAdapter() DiagnosisAdapter {
	db := driver.NewPsql()
	return &diagnosisAdapter{db: db}
}

type diagnosisAdapter struct {
	db *gorm.DB
}

func (a *diagnosisAdapter) CreateDiagnosisHistory(history *models.DiagnosisHistory) error {
	return a.db.Create(history).Error
}

func (a *diagnosisAdapter) GetDiagnosisHistoryByUserID(userID string) ([]models.DiagnosisHistory, error) {
	var histories []models.DiagnosisHistory
	err := a.db.Where("user_id = ?", userID).
		Preload("User").
		Preload("UserAvatar").
		Preload("TargetAvatar").
		Order("created_at DESC").
		Find(&histories).Error
	return histories, err
}

func (a *diagnosisAdapter) GetDiagnosisHistoryByID(id string) (models.DiagnosisHistory, error) {
	var history models.DiagnosisHistory
	err := a.db.Where("id = ?", id).
		Preload("User").
		Preload("UserAvatar").
		Preload("TargetAvatar").
		First(&history).Error
	return history, err
}

func (a *diagnosisAdapter) UpdateUserAvatarRelationPoints(userID, avatarID string, additionalPoints int) error {
	return a.db.Model(&models.UserAvatarRelation{}).
		Where("user_id = ? AND avatar_id = ?", userID, avatarID).
		Update("matching_point", gorm.Expr("matching_point + ?", additionalPoints)).Error
}

func (a *diagnosisAdapter) GetUserAvatarRelation(userID, avatarID string) (models.UserAvatarRelation, error) {
	var relation models.UserAvatarRelation
	err := a.db.Where("user_id = ? AND avatar_id = ?", userID, avatarID).First(&relation).Error
	return relation, err
}

func (a *diagnosisAdapter) CreateUserAvatarRelation(relation models.UserAvatarRelation) error {
	return a.db.Create(&relation).Error
}

func (a *diagnosisAdapter) GetAvatarByID(id string) (models.Avatar, error) {
	var avatar models.Avatar
	err := a.db.Where("id = ?", id).First(&avatar).Error
	return avatar, err
}

func (a *diagnosisAdapter) GetByUserID(userID string) (models.Avatar, error) {
	var avatar models.Avatar
	err := a.db.Where("user_id = ?", userID).First(&avatar).Error
	return avatar, err
}
