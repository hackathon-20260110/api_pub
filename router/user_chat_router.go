package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func UserChatRouter(e *echo.Echo, container *dig.Container) {
	userChatController := controller.NewUserChatController(container)

	firebaseAuth := middleware.FirebaseAuthMiddleware()

	e.GET("/user-chats/matched", userChatController.GetMatchedUsers, firebaseAuth)
	e.POST("/user-chats/:partner_id/messages", userChatController.SendMessage, firebaseAuth)
	e.GET("/user-chats/:partner_id/messages", userChatController.GetMessages, firebaseAuth)
}
