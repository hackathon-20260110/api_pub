package response

import "time"

type MatchedUser struct {
	UserID          string    `json:"user_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	DisplayName     string    `json:"display_name" example:"John Doe"`
	Gender          string    `json:"gender" example:"male"`
	Bio             string    `json:"bio" example:"Hello, I'm John!"`
	ProfileImageURL string    `json:"profile_image_url" example:"https://example.com/image.jpg"`
	MatchedAt       time.Time `json:"matched_at" example:"2024-01-01T12:00:00Z"`
}

type GetMatchedUsersResponse struct {
	MatchedUsers []MatchedUser `json:"matched_users"`
}

type UserChatMessageResponse struct {
	ID         string    `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	SenderID   string    `json:"sender_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	SenderType string    `json:"sender_type" example:"user"`
	Message    string    `json:"message" example:"Hello, nice to meet you!"`
	CreatedAt  time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
}

type SendUserChatMessageResponse struct {
	Message     string                  `json:"message" example:"Message sent successfully"`
	ChatMessage UserChatMessageResponse `json:"chat_message"`
}

type GetUserChatMessagesResponse struct {
	Messages []UserChatMessageResponse `json:"messages"`
}
