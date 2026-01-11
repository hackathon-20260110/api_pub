package response

import (
	"time"

	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/utils"
)

// LoginResponse ログインレスポンス
type LoginResponse struct {
	User User `json:"user"`
}

// User ユーザー情報
type User struct {
	ID                  string `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	DisplayName         string `json:"display_name" example:"山田太郎"`
	Age                 int    `json:"age" example:"25"`
	Gender              string `json:"gender" example:"male"`
	ProfileImageURL     string `json:"profile_image_url" example:"https://example.com/images/profile.jpg"`
	Bio                 string `json:"bio" example:"よろしくお願いします！"`
	OnboardingCompleted bool   `json:"onboarding_completed" example:"false"`
	CreatedAt           string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt           string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

func NewUserResponse(user models.User) User {
	return User{
		ID:                  user.ID,
		DisplayName:         user.DisplayName,
		Age:                 utils.CalculateAge(user.BirthDate, time.Now()),
		ProfileImageURL:     user.ProfileImageURL,
		Bio:                 user.Bio,
		OnboardingCompleted: user.IsOnboardingCompleted,
		CreatedAt:           user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           user.UpdatedAt.Format(time.RFC3339),
	}
}

// ErrorResponse エラーレスポンス
type ErrorResponse struct {
	Error   string `json:"error" example:"invalid_token"`
	Message string `json:"message" example:"Firebase IDトークンが不正または無効です"`
}
