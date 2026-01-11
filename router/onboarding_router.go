package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func OnboardingRouter(e *echo.Echo, container *dig.Container) {
	controller := controller.NewOnboardingController(container)

	// Firebase認証ミドルウェアを適用
	firebaseAuth := middleware.FirebaseAuthMiddleware()

	e.POST("/onboarding/start", controller.StartOnboardingChat, firebaseAuth)
	e.POST("/onboarding/chats/messages", controller.SendOnboardingMessage, firebaseAuth)
	e.POST("/onboarding/finish", controller.FinishOnboarding, firebaseAuth)
}
