package requests

type SendUserChatMessageRequest struct {
	Content string `json:"content" example:"Hello, nice to meet you!" binding:"required"`
}
