package main

import (
	"os"

	"github.com/hackathon-20260110/api/dicontainer"
	_ "github.com/hackathon-20260110/api/docs"
	"github.com/hackathon-20260110/api/driver"
	"github.com/hackathon-20260110/api/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Hackathon API
// @version 1.0
// @license.name MIT
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Firebase IDトークンをBearerトークンとして送信してください。例: "Bearer eyJhbGciOiJSUzI1NiIs..."
func main() {
	// ================================
	// 初期化処理
	// ================================

	// Firebase Admin SDK 初期化（失敗時は起動を止める）
	driver.NewFirebaseAuth()
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("10M"))

	container := dicontainer.GetContainer()

	// ================================
	// swaggerルート
	// ================================
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// ================================
	// APIルーティング
	// ================================
	// /debug/*
	router.DebugRouter(e, container)
	// /auth/*
	router.AuthRouter(e, container)
	// /onboarding/*
	router.OnboardingRouter(e, container)
	// /users/*
	router.UserRouter(e, container)
	// /chats/*
	router.ChatRouter(e, container)
	// /chats/{partnerUserId}/unlock-status, /chats/{partnerUserId}/match
	router.MatchRouter(e, container)
	// /matches/* (マッチング候補・マッチ後支援)
	router.MatchPostRouter(e, container)
	// /avatars/*
	router.AvatarRouter(e, container)
	// /profile/*, /users/me/profile/*, /users/{userId}/profile
	router.ProfileRouter(e, container)
	// /diagnosis/*
	router.SetupDiagnosisRouter(e, container)
	// /avatar-chats/*
	router.AvatarChatRouter(e, container)
	// /user-chats/*
	router.UserChatRouter(e, container)
	// /notification/*
	router.NotificationRouter(e, container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
