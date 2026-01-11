package adapter

import (
	"github.com/hackathon-20260110/api/driver"
	"github.com/hackathon-20260110/api/utils"

	"github.com/hackathon-20260110/api/models"
	"gorm.io/gorm"
)

type UserAdapter interface {
	GetByID(id string) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
}

func NewUserAdapter() UserAdapter {
	db := driver.NewPsql()
	return &userAdapter{db: db}
}

type userAdapter struct {
	db *gorm.DB
}

func (a *userAdapter) GetByID(id string) (models.User, error) {
	var user models.User
	if err := a.db.Where("id = ?", id).First(&user).Error; err != nil {
		return models.User{}, utils.ErrorRecordNotFound
	}
	return user, nil
}

func (a *userAdapter) Create(user models.User) (models.User, error) {
	if err := a.db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (a *userAdapter) Update(user models.User) (models.User, error) {
	if err := a.db.Save(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
