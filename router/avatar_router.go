package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func AvatarRouter(e *echo.Echo, container *dig.Container) {
	controller := controller.NewAvatarController(container)

	firebaseAuth := middleware.FirebaseAuthMiddleware()

	e.GET("/avatars", controller.GetAvatarList, firebaseAuth)
	e.POST("/avatars/:avatarId/matching-points", controller.UpdateMatchingPoint, firebaseAuth)
}
