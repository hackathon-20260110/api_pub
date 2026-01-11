package response

import "time"

type AvatarChatMessage struct {
	ID         string    `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	SenderType string    `json:"sender_type" example:"user"`
	Message    string    `json:"message" example:"Hello, how are you?"`
	CreatedAt  time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
}

type UnlockedUserInfo struct {
	ID    string `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Key   string `json:"key" example:"hobby"`
	Value string `json:"value" example:"reading"`
}

type SendAvatarChatMessageResponse struct {
	Message          string             `json:"message" example:"Message sent successfully"`
	AvatarResponse   AvatarChatMessage  `json:"avatar_response"`
	MatchingPoint    int                `json:"matching_point" example:"50"`
	PointChange      int                `json:"point_change" example:"10"`
	IsMatched        bool               `json:"is_matched" example:"false"`
	UnlockedMissions []UnlockedUserInfo `json:"unlocked_missions"`
}

type GetAvatarChatMessagesResponse struct {
	Messages      []AvatarChatMessage `json:"messages"`
	MatchingPoint int                 `json:"matching_point" example:"50"`
	IsMatched     bool                `json:"is_matched" example:"false"`
}

type GetAvatarChatStatusResponse struct {
	MatchingPoint     int                `json:"matching_point" example:"50"`
	IsMatched         bool               `json:"is_matched" example:"false"`
	UnlockedUserInfos []UnlockedUserInfo `json:"unlocked_user_infos"`
}
