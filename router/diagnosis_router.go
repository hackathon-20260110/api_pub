package router

import (
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func SetupDiagnosisRouter(e *echo.Echo, container *dig.Container) {
	var diagnosisController *controller.DiagnosisController
	container.Invoke(func(dc *controller.DiagnosisController) {
		diagnosisController = dc
	})

	// 診断API群
	diagnosisGroup := e.Group("/diagnosis")

	firebaseAuth := middleware.FirebaseAuthMiddleware()

	// 診断実行
	diagnosisGroup.POST("/execute", diagnosisController.ExecuteDiagnosis, firebaseAuth)

	// 診断履歴取得
	diagnosisGroup.GET("/history", diagnosisController.GetDiagnosisHistory, firebaseAuth)

	// 診断詳細取得
	diagnosisGroup.GET("/:diagnosisId", diagnosisController.GetDiagnosisDetail, firebaseAuth)
}
