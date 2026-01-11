package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/middleware"
	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type ProfileController struct {
	container *dig.Container
}

func NewProfileController(container *dig.Container) *ProfileController {
	return &ProfileController{container: container}
}

// @Summary 定義済みプロフィール項目一覧取得
// @Tags profile
// @Description アプリで定義されているプロフィール項目の一覧を取得
// @Security Bearer
// @Success 200 {object} response.PredefinedKeysResponse "定義済み項目一覧"
// @Failure 401 {object} response.ErrorResponse "認証エラー"
// @Router /profile/predefined-keys [get]
func (c *ProfileController) GetPredefinedKeys(ctx echo.Context) error {
	service := service.NewProfileService(c.container)
	keys := service.GetPredefinedInfoKeys()

	return ctx.JSON(http.StatusOK, &response.PredefinedKeysResponse{
		Keys: keys,
	})
}

// @Summary プロフィール項目作成
// @Tags profile
// @Description ユーザーのプロフィール項目を作成（ミッション設定可能）
// @Security Bearer
// @Param request body requests.CreateUserInfoRequest true "プロフィール項目作成リクエスト"
// @Success 200 {object} response.UserInfoResponse "作成された項目情報"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証エラー"
// @Router /users/me/profile/info [post]
func (c *ProfileController) CreateUserInfo(ctx echo.Context) error {
	userID := middleware.GetFirebaseUID(ctx)

	var req requests.CreateUserInfoRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "invalid_request",
			Message: "リクエストが不正です",
		})
	}

	service := service.NewProfileService(c.container)
	result, err := service.CreateUserInfo(ctx.Request().Context(), userID, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, result)
}

// @Summary プロフィール項目更新
// @Tags profile
// @Description ユーザーのプロフィール項目を更新
// @Security Bearer
// @Param infoId path string true "プロフィール項目ID"
// @Param request body requests.UpdateUserInfoRequest true "プロフィール項目更新リクエスト"
// @Success 200 {object} response.UserInfoResponse "更新された項目情報"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証エラー"
// @Failure 404 {object} response.ErrorResponse "プロフィール項目が見つからない"
// @Router /users/me/profile/info/{infoId} [put]
func (c *ProfileController) UpdateUserInfo(ctx echo.Context) error {
	userID := middleware.GetFirebaseUID(ctx)
	infoID := ctx.Param("infoId")

	var req requests.UpdateUserInfoRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "invalid_request",
			Message: "リクエストが不正です",
		})
	}

	service := service.NewProfileService(c.container)
	result, err := service.UpdateUserInfo(ctx.Request().Context(), userID, infoID, req)
	if err != nil {
		if err.Error() == "user info not found: record not found" {
			return ctx.JSON(http.StatusNotFound, &response.ErrorResponse{
				Error:   "not_found",
				Message: "プロフィール項目が見つかりません",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, result)
}

// @Summary プロフィール項目削除
// @Tags profile
// @Description ユーザーのプロフィール項目を削除
// @Security Bearer
// @Param infoId path string true "プロフィール項目ID"
// @Success 200 {object} map[string]string "削除成功"
// @Failure 401 {object} response.ErrorResponse "認証エラー"
// @Failure 404 {object} response.ErrorResponse "プロフィール項目が見つからない"
// @Router /users/me/profile/info/{infoId} [delete]
func (c *ProfileController) DeleteUserInfo(ctx echo.Context) error {
	userID := middleware.GetFirebaseUID(ctx)
	infoID := ctx.Param("infoId")

	service := service.NewProfileService(c.container)
	err := service.DeleteUserInfo(ctx.Request().Context(), userID, infoID)
	if err != nil {
		if err.Error() == "user info not found: record not found" {
			return ctx.JSON(http.StatusNotFound, &response.ErrorResponse{
				Error:   "not_found",
				Message: "プロフィール項目が見つかりません",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "プロフィール項目を削除しました",
	})
}

// @Summary 自分のプロフィール取得
// @Tags profile
// @Description 自分の全プロフィール情報を取得（ミッション情報含む）
// @Security Bearer
// @Success 200 {object} response.UserProfileResponse "プロフィール情報"
// @Failure 401 {object} response.ErrorResponse "認証エラー"
// @Router /users/me/profile [get]
func (c *ProfileController) GetMyProfile(ctx echo.Context) error {
	userID := middleware.GetFirebaseUID(ctx)

	service := service.NewProfileService(c.container)
	result, err := service.GetUserProfile(ctx.Request().Context(), userID, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, result)
}

// @Summary 他ユーザーのプロフィール取得
// @Tags profile
// @Description 指定ユーザーのプロフィール情報を取得（ミッション解除状況含む）
// @Security Bearer
// @Param userId path string true "対象ユーザーID"
// @Success 200 {object} response.UserProfileResponse "プロフィール情報"
// @Failure 401 {object} response.ErrorResponse "認証エラー"
// @Failure 404 {object} response.ErrorResponse "ユーザーが見つからない"
// @Router /users/{userId}/profile [get]
func (c *ProfileController) GetUserProfile(ctx echo.Context) error {
	viewerUserID := middleware.GetFirebaseUID(ctx)
	targetUserID := ctx.Param("userId")

	service := service.NewProfileService(c.container)
	result, err := service.GetUserProfile(ctx.Request().Context(), targetUserID, viewerUserID)
	if err != nil {
		if err.Error() == "user not found: record not found" {
			return ctx.JSON(http.StatusNotFound, &response.ErrorResponse{
				Error:   "not_found",
				Message: "ユーザーが見つかりません",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, result)
}
