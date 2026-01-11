package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/driver"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type DebugController struct {
	container *dig.Container
}

func NewDebugController(container *dig.Container) *DebugController {
	return &DebugController{container: container}
}

// @Summary Health check
// @Tags debug
// @Success 200 {string} string "OK! Server is healthy"
// @Router /debug/health [get]
func (c *DebugController) Health(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK! Server is healthy")
}

// @Summary Endpoints
// @Tags debug
// @Router /debug/endpoints [get]
func (c *DebugController) Endpoints(ctx echo.Context) error {
	endpoints := []map[string]string{
		{
			"name": "local_dev",
			"url":  "http://localhost:8080",
		},
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"endpoints": endpoints,
	})
}

// @Summary Echo
// @Tags debug
// @Param message body requests.DebugEchoRequest true "メッセージ"
// @Success 200 {object} response.DebugEchoResponse "OK! Echo is healthy"
// @Router /debug/echo [post]
func (c *DebugController) Echo(ctx echo.Context) error {
	req := new(requests.DebugEchoRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Failed to bind request",
		})
	}
	return ctx.JSON(http.StatusOK, &response.DebugEchoResponse{Message: req.Message})
}

// @Summary Atuch Check
// @Tags debug
// @Router /debug/atuch-check [get]
func (c *DebugController) AtuchCheck(ctx echo.Context) error {
	userID := middleware.GetFirebaseUID(ctx)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userID,
	})
}

// IDToken generates a Firebase ID token for testing purposes
// @Summary Firebase ID Tokenを生成
// @Description Firebase ID token（IDトークン）を直接生成します。このトークンはテスト・開発用途で使用できます。
// @Description このトークンはBearerトークンとしてAuthorizationヘッダーで直接使用できます。
// @Description 内部的には、カスタムトークンを生成してからFirebase REST APIで交換しています。
// @Tags debug
// @Param userId query string true "Firebase User ID (UID)"
// @Success 200 {object} response.DebugIDTokenResponse "ID token generated successfully"
// @Failure 400 {object} response.ErrorResponse "userId parameter is missing or empty"
// @Failure 500 {object} response.ErrorResponse "Failed to generate ID token"
// @Failure 503 {object} response.ErrorResponse "Firebase Web API Key not configured"
// @Router /debug/id-token [get]
func (c *DebugController) IDToken(ctx echo.Context) error {
	// Extract userId from query parameter
	userID := ctx.QueryParam("userId")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "missing_parameter",
			Message: "userIdパラメータが必要です",
		})
	}

	// Step 1: Generate custom token using Firebase Admin SDK
	customToken, err := driver.CreateCustomToken(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "custom_token_generation_failed",
			Message: "カスタムトークンの生成に失敗しました",
		})
	}

	// Step 2: Exchange custom token for ID token using Firebase REST API
	idToken, expiresIn, err := driver.ExchangeCustomTokenForIDToken(ctx.Request().Context(), customToken)
	if err != nil {
		// Check for specific error: missing API key
		if err.Error() == "firebase: FIREBASE_WEB_API_KEY environment variable is not set" {
			return ctx.JSON(http.StatusServiceUnavailable, &response.ErrorResponse{
				Error:   "api_key_not_configured",
				Message: "Firebase Web API Keyが設定されていません。管理者に連絡してください",
			})
		}

		// Generic token exchange error
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "token_exchange_failed",
			Message: "IDトークンへの交換に失敗しました",
		})
	}

	// Return success response with ID token
	return ctx.JSON(http.StatusOK, &response.DebugIDTokenResponse{
		IDToken:   idToken,
		ExpiresIn: expiresIn,
		UserID:    userID,
		TokenType: "Bearer",
	})
}
