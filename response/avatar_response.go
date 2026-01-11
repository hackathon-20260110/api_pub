package response

import "time"

type Avatar struct {
	ID                string      `json:"id"`
	UserID            string      `json:"user_id"`
	AvatarIconURL     string      `json:"avatar_icon_url"`
	Prompt            string      `json:"prompt"`
	PersonalityTraits interface{} `json:"personality_traits"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type AvatarWithRelation struct {
	Avatar          Avatar              `json:"avatar"`
	Relation        *UserAvatarRelation `json:"relation,omitempty"`
	UserDisplayName string              `json:"user_display_name"`
	UserAge         int                 `json:"user_age"`
	UserBio         string              `json:"user_bio"`
}

type UserAvatarRelation struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	AvatarID      string    `json:"avatar_id"`
	MatchingPoint int       `json:"matching_point"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AvatarListResponse struct {
	Avatars []AvatarWithRelation `json:"avatars"`
}

type UpdateMatchingPointResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	AvatarID      string    `json:"avatar_id"`
	MatchingPoint int       `json:"matching_point"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AvatarDetailResponse struct {
	Avatar Avatar `json:"avatar"`
}
