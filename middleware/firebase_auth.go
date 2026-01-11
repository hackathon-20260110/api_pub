package middleware

import (
	"net/http"
	"strings"

	"github.com/hackathon-20260110/api/driver"
	"github.com/hackathon-20260110/api/response"
	"github.com/labstack/echo/v4"
)

const (
	contextKeyFirebaseUID = "firebase_uid"
	contextKeyEmail       = "email"
)

func FirebaseAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, &response.ErrorResponse{
					Error:   "missing_token",
					Message: "Authorizationヘッダーが必要です",
				})
			}

			// Extract token from "Bearer <token>" format
			var idToken string
			if strings.HasPrefix(authHeader, "Bearer ") {
				idToken = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				// Allow token without Bearer prefix for flexibility
				idToken = authHeader
			}

			if idToken == "" {
				return c.JSON(http.StatusUnauthorized, &response.ErrorResponse{
					Error:   "empty_token",
					Message: "IDトークンが空です",
				})
			}

			// Verify the ID token using Firebase Admin SDK
			uid, email, err := driver.VerifyIDToken(c.Request().Context(), idToken)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, &response.ErrorResponse{
					Error:   "invalid_token",
					Message: "Firebase IDトークンが不正または無効です",
				})
			}

			if uid == "" {
				return c.JSON(http.StatusUnauthorized, &response.ErrorResponse{
					Error:   "invalid_token",
					Message: "Firebase IDトークンが不正または無効です",
				})
			}

			// NOTE: 既存コードの互換性のため両方セットする
			c.Set(contextKeyFirebaseUID, uid)
			c.Set("userID", uid)
			if email != "" {
				c.Set(contextKeyEmail, email)
			}
			return next(c)
		}
	}
}

// GetFirebaseUID retrieves the Firebase UID from the Echo context.
// Returns an empty string if not set.
func GetFirebaseUID(c echo.Context) string {
	if uid, ok := c.Get(contextKeyFirebaseUID).(string); ok {
		return uid
	}
	// 互換性: 旧キー（UserControllerなど）が使用
	if uid, ok := c.Get("userID").(string); ok {
		return uid
	}
	return ""
}

// GetEmail retrieves the email from the Echo context.
// Returns an empty string if not set.
func GetEmail(c echo.Context) string {
	if email, ok := c.Get(contextKeyEmail).(string); ok {
		return email
	}
	return ""
}
