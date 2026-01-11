package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func ProfileRouter(e *echo.Echo, container *dig.Container) {
	controller := controller.NewProfileController(container)
	firebaseAuth := middleware.FirebaseAuthMiddleware()

	// プロフィール関連エンドポイント
	e.GET("/profile/predefined-keys", controller.GetPredefinedKeys, firebaseAuth)
	e.POST("/users/me/profile/info", controller.CreateUserInfo, firebaseAuth)
	e.PUT("/users/me/profile/info/:infoId", controller.UpdateUserInfo, firebaseAuth)
	e.DELETE("/users/me/profile/info/:infoId", controller.DeleteUserInfo, firebaseAuth)
	e.GET("/users/me/profile", controller.GetMyProfile, firebaseAuth)
	e.GET("/users/:userId/profile", controller.GetUserProfile, firebaseAuth)
}

