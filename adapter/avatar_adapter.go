package adapter

import (
	"github.com/hackathon-20260110/api/models"
	"gorm.io/gorm"
)

type AvatarAdapter interface {
	GetByID(id string) (*models.Avatar, error)
	GetByUserID(userID string) (*models.Avatar, error)
	Create(avatar models.Avatar) (*models.Avatar, error)
	Update(avatar models.Avatar) (*models.Avatar, error)
	GetOppositeGenderAvatars(currentUserGender string) ([]models.Avatar, error)
	GetUserAvatarRelation(userID, avatarID string) (models.UserAvatarRelation, error)
	GetUserAvatarRelationsByUserID(userID string) ([]models.UserAvatarRelation, error)
	CreateUserAvatarRelation(relation models.UserAvatarRelation) error
	UpdateUserAvatarRelation(relation models.UserAvatarRelation) error
}
type avatarAdapter struct {
	db *gorm.DB
}

func NewAvatarAdapter(db *gorm.DB) AvatarAdapter {
	return &avatarAdapter{db: db}
}

func (a *avatarAdapter) GetByID(id string) (*models.Avatar, error) {
	var avatar models.Avatar
	if err := a.db.Where("id = ?", id).First(&avatar).Error; err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (a *avatarAdapter) GetByUserID(userID string) (*models.Avatar, error) {
	var avatar models.Avatar
	if err := a.db.Where("user_id = ?", userID).First(&avatar).Error; err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (a *avatarAdapter) Create(avatar models.Avatar) (*models.Avatar, error) {
	if err := a.db.Create(&avatar).Error; err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (a *avatarAdapter) Update(avatar models.Avatar) (*models.Avatar, error) {
	if err := a.db.Save(&avatar).Error; err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (a *avatarAdapter) GetOppositeGenderAvatars(currentUserGender string) ([]models.Avatar, error) {
	var avatars []models.Avatar
	var query *gorm.DB

	targetGender := map[string]string{"male": "female", "female": "male"}[currentUserGender]

	query = a.db.Joins("JOIN users ON avatars.user_id = users.id")
	if targetGender != "" {
		query = query.Where("users.gender = ?", targetGender)
	} else {
		query = query.Where("users.gender != ?", currentUserGender)
	}

	if err := query.Find(&avatars).Error; err != nil {
		return nil, err
	}
	return avatars, nil
}

func (a *avatarAdapter) GetUserAvatarRelation(userID, avatarID string) (models.UserAvatarRelation, error) {
	var relation models.UserAvatarRelation
	if err := a.db.Where("user_id = ? AND avatar_id = ?", userID, avatarID).First(&relation).Error; err != nil {
		return models.UserAvatarRelation{}, err
	}
	return relation, nil
}

func (a *avatarAdapter) GetUserAvatarRelationsByUserID(userID string) ([]models.UserAvatarRelation, error) {
	var relations []models.UserAvatarRelation
	if err := a.db.Where("user_id = ?", userID).Find(&relations).Error; err != nil {
		return nil, err
	}
	return relations, nil
}

func (a *avatarAdapter) CreateUserAvatarRelation(relation models.UserAvatarRelation) error {
	return a.db.Create(&relation).Error
}

func (a *avatarAdapter) UpdateUserAvatarRelation(relation models.UserAvatarRelation) error {
	return a.db.Save(&relation).Error
}
