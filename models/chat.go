package models

import "time"

type SenderType string

const (
	SenderTypeUser     SenderType = "user"
	SenderTypeSystem   SenderType = "system"
	SenderTypeAvatarAI SenderType = "avatar_ai"
)

type OnboardingChat struct {
	ID                    string     `json:"id" firestore:"id"`
	SenderType            SenderType `json:"sender_type" firestore:"sender_type"`
	Message               string     `json:"message" firestore:"message"`
	IsOnboardingCompleted bool       `json:"is_onboarding_completed" firestore:"is_onboarding_completed"`
	CreatedAt             time.Time  `json:"created_at" firestore:"created_at"`
}

type UserAvatarChat struct {
	ID         string     `json:"id"`
	SenderType SenderType `json:"sender_type"`
	Message    string     `json:"message"`
	Avatar     Avatar
	CreatedAt  time.Time `json:"created_at"`
}

type UserChat struct {
	ID         string     `json:"id"`
	SenderType SenderType `json:"sender_type"`
	Message    string     `json:"message"`
	User1      User
	User2      User
	CreatedAt  time.Time `json:"created_at"`
}
