package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type UserChatController struct {
	container *dig.Container
}

func NewUserChatController(container *dig.Container) *UserChatController {
	return &UserChatController{container: container}
}

// GetMatchedUsers godoc
// @Summary マッチしているユーザー一覧取得
// @Description マッチしているユーザーの一覧を取得する
// @Tags user-chat
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.GetMatchedUsersResponse "マッチしているユーザー一覧取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /user-chats/matched [get]
func (c *UserChatController) GetMatchedUsers(ctx echo.Context) error {
	userID := ctx.Get("uid").(string)

	userChatService := service.NewUserChatService(c.container)
	matchedUsers, err := userChatService.GetMatchedUsers(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to get matched users",
		})
	}

	responseUsers := make([]response.MatchedUser, 0, len(matchedUsers))
	for _, user := range matchedUsers {
		responseUsers = append(responseUsers, response.MatchedUser{
			UserID:          user.UserID,
			DisplayName:     user.DisplayName,
			Gender:          user.Gender,
			Bio:             user.Bio,
			ProfileImageURL: user.ProfileImageURL,
			MatchedAt:       user.MatchedAt,
		})
	}

	return ctx.JSON(http.StatusOK, response.GetMatchedUsersResponse{
		MatchedUsers: responseUsers,
	})
}

// SendMessage godoc
// @Summary ユーザーチャットメッセージ送信
// @Description マッチしているユーザーにメッセージを送信する
// @Tags user-chat
// @Accept json
// @Produce json
// @Security Bearer
// @Param partner_id path string true "相手ユーザーID"
// @Param request body requests.SendUserChatMessageRequest true "メッセージ内容"
// @Success 200 {object} response.SendUserChatMessageResponse "メッセージ送信成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "マッチしていないユーザーへのメッセージ送信"
// @Router /user-chats/{partner_id}/messages [post]
func (c *UserChatController) SendMessage(ctx echo.Context) error {
	userID := ctx.Get("uid").(string)
	partnerID := ctx.Param("partner_id")

	var req requests.SendUserChatMessageRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if req.Content == "" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Content is required",
		})
	}

	userChatService := service.NewUserChatService(c.container)
	message, err := userChatService.SendMessage(ctx.Request().Context(), userID, partnerID, req.Content)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, response.ErrorResponse{
			Error: "Cannot send message to this user. Make sure you are matched.",
		})
	}

	return ctx.JSON(http.StatusOK, response.SendUserChatMessageResponse{
		Message: "Message sent successfully",
		ChatMessage: response.UserChatMessageResponse{
			ID:         message.ID,
			SenderID:   message.SenderID,
			SenderType: string(message.SenderType),
			Message:    message.Message,
			CreatedAt:  message.CreatedAt,
		},
	})
}

// GetMessages godoc
// @Summary ユーザーチャットメッセージ一覧取得
// @Description マッチしているユーザーとのチャット履歴を取得する
// @Tags user-chat
// @Accept json
// @Produce json
// @Security Bearer
// @Param partner_id path string true "相手ユーザーID"
// @Success 200 {object} response.GetUserChatMessagesResponse "メッセージ一覧取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "マッチしていないユーザーのメッセージ取得"
// @Router /user-chats/{partner_id}/messages [get]
func (c *UserChatController) GetMessages(ctx echo.Context) error {
	userID := ctx.Get("uid").(string)
	partnerID := ctx.Param("partner_id")

	userChatService := service.NewUserChatService(c.container)
	messages, err := userChatService.GetMessages(ctx.Request().Context(), userID, partnerID)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, response.ErrorResponse{
			Error: "Cannot get messages with this user. Make sure you are matched.",
		})
	}

	responseMessages := make([]response.UserChatMessageResponse, 0, len(messages))
	for _, msg := range messages {
		responseMessages = append(responseMessages, response.UserChatMessageResponse{
			ID:         msg.ID,
			SenderID:   msg.SenderID,
			SenderType: string(msg.SenderType),
			Message:    msg.Message,
			CreatedAt:  msg.CreatedAt,
		})
	}

	return ctx.JSON(http.StatusOK, response.GetUserChatMessagesResponse{
		Messages: responseMessages,
	})
}
