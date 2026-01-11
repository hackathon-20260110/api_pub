package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func AuthRouter(e *echo.Echo, container *dig.Container) {
	controller := controller.NewAuthController(container)

	// Firebase認証ミドルウェアを適用
	firebaseAuth := middleware.FirebaseAuthMiddleware()

	// すべてのauthエンドポイントでFirebase認証を必須にする
	e.GET("/auth/me", controller.Me, firebaseAuth)
}
