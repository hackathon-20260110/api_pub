package controller

import (
	"fmt"
	"net/http"

	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type OnboardingController struct {
	container *dig.Container
}

func NewOnboardingController(container *dig.Container) *OnboardingController {
	return &OnboardingController{container: container}
}

// @Summary オンボーディングチャット開始
// @Tags onboarding
// @Description オンボーディングチャットを開始する
// @Security Bearer
// @Param request body requests.StartOnboardingRequest true "オンボーディングチャット開始リクエスト"
// @Success 200 {object} response.StartOnboardingResponse "オンボーディングチャット開始成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /onboarding/start [post]
func (c *OnboardingController) StartOnboardingChat(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	s := service.NewOnboardingService(c.container)
	err := s.StartOnboardingChat(ctx.Request().Context(), userID)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "オンボーディングチャット開始に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, &response.StartOnboardingResponse{
		Message: "オンボーディングチャット開始成功",
	})
}

// @Summary オンボーディングチャットメッセージ送信
// @Tags onboarding
// @Description オンボーディングチャットメッセージを送信する
// @Security Bearer
// @Param request body requests.SendOnboardingMessageRequest true "オンボーディングチャットメッセージ送信リクエスト"
// @Success 200 {object} response.SendOnboardingMessageResponse "オンボーディングチャットメッセージ送信成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /onboarding/chats/messages [post]
func (c *OnboardingController) SendOnboardingMessage(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	var req requests.SendOnboardingMessageRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "bad_request",
			Message: "リクエストが不正です",
		})
	}

	s := service.NewOnboardingService(c.container)
	err := s.SendOnboardingMessage(ctx.Request().Context(), userID, req.Content)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "オンボーディングチャットメッセージ送信に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, &response.SendOnboardingMessageResponse{
		Message: "メッセージを受け付けました",
	})
}

// @Summary オンボーディングチャット完了
// @Tags onboarding
// @Description オンボーディングチャットを完了する
// @Security Bearer
// @Success 200 {object} response.OnboardingCompleteResponse "オンボーディングチャット完了成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /onboarding/finish [post]
func (c *OnboardingController) FinishOnboarding(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	s := service.NewOnboardingService(c.container)
	userInfos, avatar, err := s.FinishOnboarding(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "オンボーディングチャット完了に失敗しました",
		})
	}

	userInfoValues := make([]models.UserInfo, 0, len(userInfos))
	for _, ui := range userInfos {
		if ui != nil {
			userInfoValues = append(userInfoValues, *ui)
		}
	}

	return ctx.JSON(http.StatusOK, &response.OnboardingCompleteResponse{
		UserInfos: userInfoValues,
		Avatar:    *avatar,
	})
}
