package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type AuthController struct {
	container *dig.Container
}

func NewAuthController(container *dig.Container) *AuthController {
	return &AuthController{container: container}
}

// @Summary ログイン中ユーザー取得
// @Tags auth
// @Description AuthorizationヘッダーのFirebase IDトークンからログイン中のユーザー情報を取得する
// @Security Bearer
// @Success 200 {object} response.User "OK"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /auth/me [get]
func (c *AuthController) Me(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDとEmailを取得
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	s := service.NewUserService(c.container)
	u, err := s.GetUserByID(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "ユーザー情報の取得に失敗しました",
		})
	}
	return ctx.JSON(http.StatusOK, u)

}
