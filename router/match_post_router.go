package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func MatchPostRouter(e *echo.Echo, container *dig.Container) {
	controller := controller.NewMatchPostController(container)

	// Firebase認証ミドルウェアを適用
	firebaseAuth := middleware.FirebaseAuthMiddleware()

	// マッチング関連のエンドポイント
	e.GET("/matches/candidates", controller.GetCandidates, firebaseAuth) // 相手候補一覧（/partnersから移動）
	e.GET("/matches", controller.GetMatches, firebaseAuth)               // マッチ一覧
	e.GET("/matches/:id/suggestions", controller.GetSuggestions, firebaseAuth)
	e.POST("/matches/:id/messages", controller.SendMatchMessage, firebaseAuth)
	e.GET("/matches/:id/messages", controller.GetMatchMessages, firebaseAuth)
	e.GET("/matches/:id/reply-assist", controller.ReplyAssist, firebaseAuth)
}
