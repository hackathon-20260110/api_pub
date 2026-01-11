package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type NotificationController struct {
	container *dig.Container
}

func NewNotificationController(container *dig.Container) *NotificationController {
	return &NotificationController{container: container}
}

// MarkAsRead godoc
// @Summary 通知を既読にする
// @Description 指定された通知IDの通知を既読状態にします
// @Tags notification
// @Accept json
// @Produce json
// @Param id path string true "通知ID"
// @Success 200 {object} response.MarkAsReadResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security Bearer
// @Router /notification/read/{id} [post]
func (c *NotificationController) MarkAsRead(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	notificationID := ctx.Param("id")
	if notificationID == "" {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "bad_request",
			Message: "通知IDが必要です",
		})
	}

	notificationService := service.NewNotificationService(c.container)
	if err := notificationService.MarkAsRead(ctx.Request().Context(), userID, notificationID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "通知の既読処理に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, &response.MarkAsReadResponse{
		Message: "通知を既読にしました",
	})
}
