package requests

type SendAvatarChatMessageRequest struct {
	Content string `json:"content" example:"Hello, how are you?" binding:"required"`
}
