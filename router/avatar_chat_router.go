package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func AvatarChatRouter(e *echo.Echo, container *dig.Container) {
	c := controller.NewAvatarChatController(container)

	firebaseAuth := middleware.FirebaseAuthMiddleware()

	e.POST("/avatar-chats/:avatar_id/messages", c.SendMessage, firebaseAuth)
	e.GET("/avatar-chats/:avatar_id/messages", c.GetMessages, firebaseAuth)
	e.GET("/avatar-chats/:avatar_id/status", c.GetStatus, firebaseAuth)
}
