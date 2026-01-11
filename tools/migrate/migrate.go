package main

import (
	"github.com/hackathon-20260110/api/driver"
	"github.com/hackathon-20260110/api/models"
)

func main() {
	db := driver.NewPsql()

	db.AutoMigrate(&models.Avatar{})
	db.AutoMigrate(&models.UserInfo{})
	db.AutoMigrate(&models.Mission{})
	db.AutoMigrate(&models.MissionUnlock{})
	db.AutoMigrate(&models.UserAvatarRelation{})
	db.AutoMigrate(&models.Matching{})
	db.AutoMigrate(&models.DiagnosisHistory{})
	db.AutoMigrate(&models.User{})
}
