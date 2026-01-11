package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func DebugRouter(e *echo.Echo, container *dig.Container) {
	controller := controller.NewDebugController(container)
	firebaseAuth := middleware.FirebaseAuthMiddleware()

	e.GET("/debug/health", controller.Health)
	e.GET("/debug/endpoints", controller.Endpoints)
	e.POST("/debug/echo", controller.Echo)
	e.GET("/debug/auth-check", controller.AtuchCheck, firebaseAuth)
	e.GET("/debug/id-token", controller.IDToken)
}
