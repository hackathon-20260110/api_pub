package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func MatchRouter(e *echo.Echo, container *dig.Container) {
	controller := controller.NewMatchController(container)

	// Firebase認証ミドルウェアを適用
	firebaseAuth := middleware.FirebaseAuthMiddleware()

	// マッチング関連のエンドポイント（チャットのサブパスとして定義、IDは相手のUserID）
	e.GET("/chats/:partnerUserId/unlock-status", controller.GetUnlockStatus, firebaseAuth)
	e.POST("/chats/:partnerUserId/match", controller.Match, firebaseAuth)
}
