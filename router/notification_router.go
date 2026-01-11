package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func NotificationRouter(e *echo.Echo, container *dig.Container) {
	notificationController := controller.NewNotificationController(container)
	firebaseAuth := middleware.FirebaseAuthMiddleware()

	e.POST("/notification/read/:id", notificationController.MarkAsRead, firebaseAuth)
}
