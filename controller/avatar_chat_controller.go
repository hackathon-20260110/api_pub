package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type AvatarChatController struct {
	container *dig.Container
}

func NewAvatarChatController(container *dig.Container) *AvatarChatController {
	return &AvatarChatController{container: container}
}

// @Summary アバターチャットメッセージ送信
// @Tags avatar-chat
// @Description アバターにメッセージを送信し、応答を受け取る
// @Security Bearer
// @Param avatar_id path string true "アバターID"
// @Param request body requests.SendAvatarChatMessageRequest true "メッセージ送信リクエスト"
// @Success 200 {object} response.SendAvatarChatMessageResponse "メッセージ送信成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /avatar-chats/{avatar_id}/messages [post]
func (c *AvatarChatController) SendMessage(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	avatarID := ctx.Param("avatar_id")
	if avatarID == "" {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "bad_request",
			Message: "アバターIDが必要です",
		})
	}

	var req requests.SendAvatarChatMessageRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "bad_request",
			Message: "リクエストが不正です",
		})
	}

	if req.Content == "" {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "bad_request",
			Message: "メッセージ内容が必要です",
		})
	}

	s := service.NewAvatarChatService(c.container)
	result, err := s.SendMessage(ctx.Request().Context(), userID, avatarID, req.Content)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "メッセージ送信に失敗しました",
		})
	}

	unlockedMissions := make([]response.UnlockedUserInfo, 0, len(result.UnlockedMissions))
	for _, m := range result.UnlockedMissions {
		unlockedMissions = append(unlockedMissions, response.UnlockedUserInfo{
			ID:    m.UserInfoID,
			Key:   m.Key,
			Value: m.Value,
		})
	}

	return ctx.JSON(http.StatusOK, &response.SendAvatarChatMessageResponse{
		Message: "メッセージを送信しました",
		AvatarResponse: response.AvatarChatMessage{
			ID:         result.AvatarResponse.ID,
			SenderType: string(result.AvatarResponse.SenderType),
			Message:    result.AvatarResponse.Message,
			CreatedAt:  result.AvatarResponse.CreatedAt,
		},
		MatchingPoint:    result.MatchingPoint,
		PointChange:      result.PointChange,
		IsMatched:        result.IsMatched,
		UnlockedMissions: unlockedMissions,
	})
}

// @Summary アバターチャットメッセージ取得
// @Tags avatar-chat
// @Description アバターとのチャット履歴を取得する
// @Security Bearer
// @Param avatar_id path string true "アバターID"
// @Success 200 {object} response.GetAvatarChatMessagesResponse "チャット履歴取得成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /avatar-chats/{avatar_id}/messages [get]
func (c *AvatarChatController) GetMessages(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	avatarID := ctx.Param("avatar_id")
	if avatarID == "" {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "bad_request",
			Message: "アバターIDが必要です",
		})
	}

	s := service.NewAvatarChatService(c.container)
	messages, matchingPoint, isMatched, err := s.GetMessages(ctx.Request().Context(), userID, avatarID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "チャット履歴取得に失敗しました",
		})
	}

	responseMessages := make([]response.AvatarChatMessage, 0, len(messages))
	for _, m := range messages {
		responseMessages = append(responseMessages, response.AvatarChatMessage{
			ID:         m.ID,
			SenderType: string(m.SenderType),
			Message:    m.Message,
			CreatedAt:  m.CreatedAt,
		})
	}

	return ctx.JSON(http.StatusOK, &response.GetAvatarChatMessagesResponse{
		Messages:      responseMessages,
		MatchingPoint: matchingPoint,
		IsMatched:     isMatched,
	})
}

// @Summary アバターチャットステータス取得
// @Tags avatar-chat
// @Description アバターとのマッチングステータスを取得する
// @Security Bearer
// @Param avatar_id path string true "アバターID"
// @Success 200 {object} response.GetAvatarChatStatusResponse "ステータス取得成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /avatar-chats/{avatar_id}/status [get]
func (c *AvatarChatController) GetStatus(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	avatarID := ctx.Param("avatar_id")
	if avatarID == "" {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "bad_request",
			Message: "アバターIDが必要です",
		})
	}

	s := service.NewAvatarChatService(c.container)
	matchingPoint, isMatched, unlockedMissions, err := s.GetStatus(ctx.Request().Context(), userID, avatarID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "ステータス取得に失敗しました",
		})
	}

	unlockedUserInfos := make([]response.UnlockedUserInfo, 0, len(unlockedMissions))
	for _, m := range unlockedMissions {
		unlockedUserInfos = append(unlockedUserInfos, response.UnlockedUserInfo{
			ID:    m.UserInfoID,
			Key:   m.Key,
			Value: m.Value,
		})
	}

	return ctx.JSON(http.StatusOK, &response.GetAvatarChatStatusResponse{
		MatchingPoint:     matchingPoint,
		IsMatched:         isMatched,
		UnlockedUserInfos: unlockedUserInfos,
	})
}
