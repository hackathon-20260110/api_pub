package adapter

import (
	"github.com/hackathon-20260110/api/models"
	"gorm.io/gorm"
)

type MatchingAdapter interface {
	GetMatchingByUsers(user1ID string, user2ID string) (*models.Matching, error)
	CreateMatching(matching models.Matching) error
	GetMatchingsByUserID(userID string) ([]models.Matching, error)
}

type matchingAdapter struct {
	db *gorm.DB
}

func NewMatchingAdapter(db *gorm.DB) MatchingAdapter {
	return &matchingAdapter{db: db}
}

func (a *matchingAdapter) GetMatchingByUsers(user1ID string, user2ID string) (*models.Matching, error) {
	var matching models.Matching
	if err := a.db.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", user1ID, user2ID, user2ID, user1ID).First(&matching).Error; err != nil {
		return nil, err
	}
	return &matching, nil
}

func (a *matchingAdapter) CreateMatching(matching models.Matching) error {
	return a.db.Create(&matching).Error
}

func (a *matchingAdapter) GetMatchingsByUserID(userID string) ([]models.Matching, error) {
	var matchings []models.Matching
	if err := a.db.Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&matchings).Error; err != nil {
		return nil, err
	}
	return matchings, nil
}
