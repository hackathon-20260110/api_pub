package adapter

import (
	"github.com/hackathon-20260110/api/models"
	"gorm.io/gorm"
)

type UserInfoAdapter interface {
	GetByID(id string) (*models.UserInfo, error)
	GetByUserID(userID string) ([]*models.UserInfo, error)
	Create(userInfo models.UserInfo) (*models.UserInfo, error)
	CreateMany(userInfos []*models.UserInfo) ([]*models.UserInfo, error)
	Update(userInfo models.UserInfo) (*models.UserInfo, error)
}

type userInfoAdapter struct {
	db *gorm.DB
}

func NewUserInfoAdapter(db *gorm.DB) UserInfoAdapter {
	return &userInfoAdapter{db: db}
}

func (a *userInfoAdapter) GetByID(id string) (*models.UserInfo, error) {
	var userInfo models.UserInfo
	if err := a.db.Where("id = ?", id).First(&userInfo).Error; err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (a *userInfoAdapter) GetByUserID(userID string) ([]*models.UserInfo, error) {
	var userInfos []*models.UserInfo
	if err := a.db.Where("user_id = ?", userID).Find(&userInfos).Error; err != nil {
		return nil, err
	}
	return userInfos, nil
}

func (a *userInfoAdapter) Create(userInfo models.UserInfo) (*models.UserInfo, error) {
	if err := a.db.Create(&userInfo).Error; err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (a *userInfoAdapter) CreateMany(userInfos []*models.UserInfo) ([]*models.UserInfo, error) {
	if err := a.db.Create(&userInfos).Error; err != nil {
		return nil, err
	}
	return userInfos, nil
}

func (a *userInfoAdapter) Update(userInfo models.UserInfo) (*models.UserInfo, error) {
	if err := a.db.Save(&userInfo).Error; err != nil {
		return nil, err
	}
	return &userInfo, nil
}
