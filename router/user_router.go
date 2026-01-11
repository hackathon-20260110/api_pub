package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func UserRouter(e *echo.Echo, container *dig.Container) {
	controller := controller.NewUserController(container)

	// Firebase認証ミドルウェアを適用
	firebaseAuth := middleware.FirebaseAuthMiddleware()

	// すべてのusersエンドポイントでFirebase認証を必須にする
	e.GET("/users/me", controller.GetMe, firebaseAuth)
	e.POST("/users", controller.CreateUser, firebaseAuth)
	e.PUT("/users/me", controller.UpdateMe, firebaseAuth)
	e.GET("/users/:userId", controller.GetUser, firebaseAuth)
	e.GET("/users/:userId/avatar-ai", controller.GetUserAvatarAI, firebaseAuth)
	e.GET("/users/me/avatar-ai", controller.GetMyAvatarAI, firebaseAuth)
}
