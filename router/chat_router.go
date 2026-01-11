package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func ChatRouter(e *echo.Echo, container *dig.Container) {
	controller := controller.NewChatController(container)

	// Firebase認証ミドルウェアを適用
	firebaseAuth := middleware.FirebaseAuthMiddleware()

	// すべてのchatsエンドポイントでFirebase認証を必須にする
	e.POST("/chats", controller.CreateChat, firebaseAuth)
	e.GET("/chats", controller.GetChats, firebaseAuth)
	e.GET("/chats/:partnerUserId", controller.GetChatDetail, firebaseAuth)
	e.POST("/chats/:partnerUserId/messages", controller.SendMessage, firebaseAuth)
	e.GET("/chats/:partnerUserId/score", controller.GetChatScore, firebaseAuth)
	// メッセージ履歴はFirestoreから直接取得するため、APIエンドポイントは不要
}
