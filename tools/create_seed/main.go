package main

import (
	"time"

	"github.com/hackathon-20260110/api/driver"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/utils"
	"gorm.io/gorm"
)

var uid1 = utils.GenerateULID()
var uid2 = utils.GenerateULID()
var uid3 = utils.GenerateULID()
var uid4 = utils.GenerateULID()
var uid5 = utils.GenerateULID()
var uid6 = utils.GenerateULID()

func main() {
	db := driver.NewPsql()

	createUsers(db)
	createUserInfos(db)
	createAvatars(db)
	createUserInfoWithMission(db)
}

func createUsers(db *gorm.DB) {
	users := []models.User{}
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	users = append(users, models.User{
		ID:                    uid1,
		DisplayName:           "はるか",
		Gender:                "female",
		BirthDate:             birthday,
		Bio:                   "よろしくお願いします！",
		ProfileImageURL:       "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_1wvs931wvs931wvs.png",
		IsOnboardingCompleted: true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	})

	users = append(users, models.User{
		ID:                    uid2,
		DisplayName:           "えりか",
		Gender:                "female",
		BirthDate:             birthday,
		Bio:                   "よろしくお願いします！",
		ProfileImageURL:       "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_f3d5mzf3d5mzf3d5.png",
		IsOnboardingCompleted: true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	})

	users = append(users, models.User{
		ID:                    uid3,
		DisplayName:           "みなみ",
		Gender:                "female",
		BirthDate:             birthday,
		Bio:                   "よろしくお願いします！",
		ProfileImageURL:       "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_is3czis3czis3czi.png",
		IsOnboardingCompleted: true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	})

	users = append(users, models.User{
		ID:                    uid4,
		DisplayName:           "ひなこ",
		Gender:                "female",
		BirthDate:             birthday,
		Bio:                   "よろしくお願いします！",
		ProfileImageURL:       "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_jvnl1tjvnl1tjvnl.png",
		IsOnboardingCompleted: true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	})

	users = append(users, models.User{
		ID:                    uid5,
		DisplayName:           "めい",
		Gender:                "female",
		BirthDate:             birthday,
		Bio:                   "よろしくお願いします！",
		ProfileImageURL:       "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_stbboistbboistbb.png",
		IsOnboardingCompleted: true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	})

	users = append(users, models.User{
		ID:                    uid6,
		DisplayName:           "のぞみ",
		Gender:                "female",
		BirthDate:             birthday,
		Bio:                   "よろしくお願いします！",
		ProfileImageURL:       "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_vih3fbvih3fbvih3.png",
		IsOnboardingCompleted: true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	})

	for _, user := range users {
		db.Create(&user)
	}
}

func createAvatars(db *gorm.DB) {
	avatars := []models.Avatar{}
	avatar1 := models.Avatar{
		ID:                utils.GenerateULID(),
		UserID:            uid1,
		AvatarIconURL:     "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_1wvs931wvs931wvs.png",
		Prompt:            "",
		PersonalityTraits: "{}",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	avatars = append(avatars, avatar1)

	avatar2 := models.Avatar{
		ID:                utils.GenerateULID(),
		UserID:            uid2,
		AvatarIconURL:     "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_f3d5mzf3d5mzf3d5.png",
		Prompt:            "",
		PersonalityTraits: "{}",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	avatars = append(avatars, avatar2)

	avatar3 := models.Avatar{
		ID:                utils.GenerateULID(),
		UserID:            uid3,
		AvatarIconURL:     "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_is3czis3czis3czi.png",
		Prompt:            "",
		PersonalityTraits: "{}",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	avatars = append(avatars, avatar3)

	avatar4 := models.Avatar{
		ID:                utils.GenerateULID(),
		UserID:            uid4,
		AvatarIconURL:     "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_jvnl1tjvnl1tjvnl.png",
		Prompt:            "",
		PersonalityTraits: "{}",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	avatars = append(avatars, avatar4)

	avatar5 := models.Avatar{
		ID:                utils.GenerateULID(),
		UserID:            uid5,
		AvatarIconURL:     "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_stbboistbboistbb.png",
		Prompt:            "",
		PersonalityTraits: "{}",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	avatars = append(avatars, avatar5)

	avatar6 := models.Avatar{
		ID:                utils.GenerateULID(),
		UserID:            uid6,
		AvatarIconURL:     "https://pub-76456fd842e04babbf3d76b005b281d5.r2.dev/Gemini_Generated_Image_vih3fbvih3fbvih3.png",
		Prompt:            "",
		PersonalityTraits: "{}",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	avatars = append(avatars, avatar6)
	for _, avatar := range avatars {
		db.Create(&avatar)
	}
}

func createUserInfos(db *gorm.DB) {
	userInfos := []models.UserInfo{}

	info1 := models.UserInfo{
		ID:              "",
		UserID:          "",
		InfoType:        models.UserInfoTypeText,
		Key:             "hobby",
		Value:           "ゲーム",
		IsMissionReward: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	userInfos = append(userInfos, info1)

	info2 := models.UserInfo{
		ID:              "",
		UserID:          "",
		InfoType:        models.UserInfoTypeText,
		Key:             "skill",
		Value:           "料理",
		IsMissionReward: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	userInfos = append(userInfos, info2)

	info3 := models.UserInfo{
		ID:              "",
		UserID:          "",
		InfoType:        models.UserInfoTypeText,
		Key:             "血液型",
		Value:           "A型",
		IsMissionReward: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	userInfos = append(userInfos, info3)

	info4 := models.UserInfo{
		ID:              "",
		UserID:          "",
		InfoType:        models.UserInfoTypeText,
		Key:             "好きな食べ物",
		Value:           "イタリアン",
		IsMissionReward: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	userInfos = append(userInfos, info4)

	uids := []string{uid1, uid2, uid3, uid4, uid5, uid6}
	for _, info := range userInfos {
		for _, uid := range uids {
			info.ID = utils.GenerateULID()
			info.UserID = uid
			db.Create(&info)
		}
	}
}

func createUserInfoWithMission(db *gorm.DB) {
	info := models.UserInfo{
		ID:              utils.GenerateULID(),
		UserID:          uid1,
		InfoType:        models.UserInfoTypeText,
		Key:             "彼氏いない歴",
		Value:           "3年",
		IsMissionReward: true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	db.Create(&info)

	thresholdPointCondition := 100
	mission := models.Mission{
		ID:                      utils.GenerateULID(),
		MissionOwnerUserID:      uid1,
		UserInfoID:              info.ID,
		ThresholdPointCondition: &thresholdPointCondition,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}
	db.Create(&mission)

}
